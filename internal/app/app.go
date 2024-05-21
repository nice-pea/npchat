package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/saime-0/cute-chat-backend/internal/config"
	"github.com/saime-0/cute-chat-backend/internal/repository/postgres"
	"sync"
)

func Start(ctx context.Context, wg sync.WaitGroup, cfg *config.Config) error {
	db, err := pgx.Connect(ctx, cfg.DbConnString())
	if err != nil {
		return fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	wg.Add(1)
	go func() {
		<-ctx.Done()
		db.Close(context.Background())
		wg.Done()
	}()

	bus := eventbus.New()
	//trans := transaction.New[*sqlx.Tx]
	//customerWrite := postgres.NewCustomerRepository(trans)

	commonRepository := postgres.NewCommonRepository(db)

	updateCustomer := command.NewUpdateCustomer(customerWrite, commonRepository)
	allCustomers := query.NewAllCustomers(commonRepository)
	customerController := web.NewCustomerController(
		updateCustomer,
		allCustomers,
	)

	registrationWrite := postgres.NewRegistrationRepository(trans)
	createRegistration := command.NewCreateRegistration(registrationWrite, commonRepository)
	confirmRegistration := command.NewConfirmRegistration(registrationWrite)

	bus.Subscribe(registration.EventEmailVerified, event.NewEmailVerifiedHandler(customerWrite, commonRepository).Handle)

	registrationController := web.NewRegistrationController(createRegistration, confirmRegistration)
	StartWebServer(ctx, lock, cfg.WebConfig, customerController, registrationController)

	emailClient := fake.NewEmailClient()
	emailGateway := fakeemail.NewClient(emailClient)
	sendEmail := command.NewSendEmail("http://localhost"+cfg.Port+"/registrations/", emailGateway)
	registrationHandler := fakesub.NewRegistrationController(sendEmail)
	mq := StartMQ(ctx, lock, registrationHandler)

	pub := fakepub.NewFakePublisher(mq)
	outboxMan := outbox.New(trans, 5, pub)
	bus.Subscribe(registration.EventRegistrationCreated, func(ctx context.Context, de eventbus.DomainEvent) error {
		// DEMO: transform the incoming domain event into an integration event if there is a need to.
		// In this case there is no need to.
		return outboxMan.Create(ctx, de)
	})
	outboxMan.Start(ctx, lock, cfg.OutboxHeartbeat)
}
