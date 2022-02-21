build/:
	mkdir build

build/libbridge.so: ../*.go *.go build/
	go build -buildmode=c-shared -o $@

build/__init__.py: *.py build/
	cp build_bridge.py build
	cd build
	./build_bridge.py
	cd ..
	cp bridge_init.py build/__init__.py

.venv/:
	python3 -m venv .venv 

.venv/touchfile: .venv/
	. .venv/bin/activate;
	pip3 install -r requirements.txt;
	touch .venv/touchfile


.PHONY: check
check: .venv/touchfile build/libbridge.so build/__init__.py
	. .venv/bin/activate; pytest --hypothesis-show-statistics

.PHONY: clean
clean:
	rm -rf _bridge.py
	rm -rf build/
	rm -rf __pycache__/
	rm -rf .hypothesis
