.PHONY: clean compile jar test

JAVA_SRC = TestApp.java
PACKAGE_DIR = testdata
CLASS_FILES = $(PACKAGE_DIR)/TestApp.class
JAR_FILE = testapp.jar
MAIN_CLASS = testdata.TestApp

# Ensure Java 8 compatibility
JAVAC_FLAGS = -source 8 -target 8

compile: $(CLASS_FILES)

$(CLASS_FILES): $(JAVA_SRC)
	@echo "Compiling Java source with Java 8 compatibility..."
	@mkdir -p $(PACKAGE_DIR)
	javac $(JAVAC_FLAGS) -d . $(JAVA_SRC)

jar: $(JAR_FILE)

$(JAR_FILE): $(CLASS_FILES)
	@echo "Creating JAR file..."
	@echo "Manifest-Version: 1.0" > manifest.txt
	@echo "Main-Class: $(MAIN_CLASS)" >> manifest.txt
	@echo "" >> manifest.txt
	jar cfm $(JAR_FILE) manifest.txt $(PACKAGE_DIR)/*.class
	@rm manifest.txt

test: jar
	@echo "Testing JAR file..."
	java -jar $(JAR_FILE) --test arg1 arg2

clean:
	@echo "Cleaning up..."
	rm -rf $(PACKAGE_DIR) $(JAR_FILE) manifest.txt

# Check if we have Java 8
check-java:
	@echo "Checking Java version..."
	@java -version 2>&1 | head -n 1
	@javac -version 2>&1