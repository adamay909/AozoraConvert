/*
azrconvert converts texts provided by [Aozora Bunko] to html, epub, and azw3 formatted for vertical typesetting.

# Usage

The text to be converted is given by the URL of the xhtml text on Aozora Bunko's site. The conversion target is specified through flags. The available targets are:

	-epub
		Produces an epub3 file.

	-kindle
			Produces an azw3 file suitable for Kindle e-readers.

	-web
			Produces a zip file containing an html file and
			all files necessary to display the page as
			intended. The difference to the original is that
			the encoding is UTF-8, and that various comments
			are removed and/or replaced to make it easier
			to control the appearance through css styling.

For example, you can convert Akutagawa Ryunosuke's Imogayu to epub by

	$ azrconvert -epub https://www.aozora.gr.jp/cards/000879/files/55_14824.html

You will get the output:

	Output written to 芋粥.epub.

You can also specify multiple targets at once. E.g.:

	$ azrconvert -epub -kindle  https://www.aozora.gr.jp/cards/000879/files/55_14824.html

	Output written to 芋粥.epub.
	Output written to 芋粥.azw3.

	$

The output file name defaults to title followed by the appropriate extension if the title can be determined from the xhtml file. If the title cannot be found, it uses the basename of the url followed by the appropriate extension.

You can manually specify the output file name using the flag

	-of
		Name of the output file.

E.g.,

	$ azrconvert -web -o imogayu  https://www.aozora.gr.jp/cards/000879/files/55_14824.html

	Output written to imogayu.zip.

You can in any case see the output file name on the command line.

# Notes

Modern web browsers have no difficulty displaying Japanese vertically but the choice of font can matter. If the display looks weird, change the serif font for Japanese to something different. For example, Noto Serif JP, Noto Sans JP, IPA fonts, work well.

[Aozora Bunko]: https://www.aozora.gr.jp
*/
package main
