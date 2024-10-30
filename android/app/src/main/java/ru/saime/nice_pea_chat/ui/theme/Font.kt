package ru.saime.nice_pea_chat.ui.theme

import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.sp

object Font {
    private val S12W500 = TextStyle(
        fontSize = 12.sp,
        fontWeight = FontWeight.W500,
    )
    private val S12W400 = TextStyle(
        fontSize = 12.sp,
        fontWeight = FontWeight.W400,
    )
    private val S14W400 = TextStyle(
        fontSize = 14.sp,
        fontWeight = FontWeight.W400,
    )
    private val S14W500 = TextStyle(
        fontSize = 14.sp,
        fontWeight = FontWeight.W500,
    )
    private val S16W400 = TextStyle(
        fontSize = 16.sp,
        fontWeight = FontWeight.W400,
    )
    private val S18W400 = TextStyle(
        fontSize = 18.sp,
        fontWeight = FontWeight.W400,
    )
    private val S20W400 = TextStyle(
        fontSize = 20.sp,
        fontWeight = FontWeight.W400,
    )

    // White
    val White12W500 = S12W500.copy(White)
    val White12W400 = S12W400.copy(White)
    val White14W400 = S14W400.copy(White)
    val White16W400 = S16W400.copy(White)
    val White16W400Italic = S16W400.copy(White, fontStyle = FontStyle.Italic)
    val White18W400 = S18W400.copy(White)
    val White20W400 = S20W400.copy(White)

    // GrayCharcoal
    val GrayCharcoal12W400 = S12W400.copy(GrayCharcoal)
    val GrayCharcoal12W400Italic = S12W400.copy(GrayCharcoal, fontStyle = FontStyle.Italic)
    val GrayCharcoal16W400 = S16W400.copy(GrayCharcoal)

    // Pink
    val Pink14W500Italic = S14W500.copy(Pink, fontStyle = FontStyle.Italic)

    // BlueSky
    val BlueSky14W400 = S14W400.copy(BlueSky)

    // BlueGraph
    val BlueGraph12W400 = S12W400.copy(BlueGraph)
    val BlueGraph14W500 = S14W500.copy(BlueGraph)

    // DarkGraph
    val DarkGraph14W500Italic = S14W500.copy(DarkGraph, fontStyle = FontStyle.Italic)
    val DarkGraph12W400 = S12W400.copy(DarkGraph)
}