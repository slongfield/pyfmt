build/:
	mkdir build

build/libbridge.so: *.go build/
	go build -buildmode=c-shared -o $@

build/__init__.py: *.py build/
	cp build_bridge.py build
	cd build
	./build_bridge.py
	cd ..
	cp bridge_init.py build/__init__.py

.PHONY: check
check: build/libbridge.so build/__init__.py
	pytest

.PHONY: clean
clean:
	rm -rf _bridge.py
	rm -rf build/
	rm -rf __pycache__/
