package build

/*
seek(tamaño, posición)
si tamaño = 0 <--- se sobrescriba
*/

var TamañoMBR int = 159
var TamañoEbr int = 30

var Alfabeto = [26]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
	'Y', 'Z'}

var Letra int = 0

var NombreArchivo string = "MIA/P1/" + string(Alfabeto[Letra]) + ".dsk"
