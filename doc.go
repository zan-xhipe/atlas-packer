/*
atlas-packer [options] spriteDir output

	--dimensions=1024x1024: Size of final atlas
	--space=1: Space inserted between sprites, in pixels
    --force=false: always rebuild atlas
    -v=false: Verbose output

outputs two files: <output>.png, and <output>.json

By default it will only build the atlas if necessary

spriteDir can only contain sprites, anything else will cause an error.
*/
package main
