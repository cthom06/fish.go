all:
	make -C lib install
	make -C cmd all
install:
	make -C lib install
	make -C cmd install
clean:
	make -C lib clean
	make -C cmd clean
nuke:
	make -C lib nuke
	make -C cmd nuke
