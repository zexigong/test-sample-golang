import json
import datetime
import argparse
import os


def find_files(repo_path: str):
    """
    Automatically finds test, source, and dependency files based on the given repository structure.
    Matches each test file with its corresponding source and dependency files.
    """
    test_source_map = {}

    # Walk through the repo and find files
    for root, _, files in os.walk(repo_path):
        for file in files:
            file_path = os.path.join(root, file)

            # Identify test files (inside any "test_*" folder)
            if "test_" in os.path.basename(root):
                test_folder = root
                test_file = file_path

                # Locate corresponding source and dependency folders
                source_folder = os.path.join(test_folder, "source_files")
                dependent_folder = os.path.join(test_folder, "dependent_files")

                # Find source files
                source_files = [
                    os.path.join(source_folder, f) for f in os.listdir(source_folder)
                    if f.endswith(".py")
                ] if os.path.exists(source_folder) else []

                # Find dependency files (or use empty.txt if none found)
                dependency_files = [
                    os.path.join(dependent_folder, f) for f in os.listdir(dependent_folder)
                    if f.endswith(".py")
                ] if os.path.exists(dependent_folder) else ["empty.txt"]

                # Store in the mapping
                test_source_map[test_file] = {
                    "sources": source_files,
                    "dependencies": dependency_files
                }

    return test_source_map


def read_files(file_paths):
    """
    Reads and returns content from multiple files.
    """
    contents = []
    for path in file_paths:
        if path == "empty.txt":
            contents.append("")  # Empty string for missing dependencies
        else:
            with open(path, "r", encoding="utf-8") as f:
                contents.append(f.read())
    return contents


def generate_messages(repository: str,
                      source_file_contents: list,
                      test_file_path: str,
                      language: str,
                      framework: str,
                      dependencies_file_names: list,
                      dependencies_file_contents: list,
                      test_example_content: str) -> dict:
    """
    Generate a JSON object in conversation format for fine-tuning.
    """
    current_time = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    system_message = (
        "You are an AI agent expert in writing unit tests. "
        "Your task is to write unit tests for the given code files of the repository. "
        "Make sure the tests can be executed without lint or compile errors."
    )

    # Combine multiple source files
    source_content = "\n\n".join(
        [f"### Source File Content:\n{content}" for content in source_file_contents]
    )

    # Combine dependencies using file names
    dependencies_content = "\n\n".join(
        [f"### Dependency File: {os.path.basename(name)}\n{content}" for name, content in zip(dependencies_file_names, dependencies_file_contents)]
    )

    user_message = (
        "### Task Information\n"
        "Based on the source code, write/rewrite tests to cover the source code.\n"
        f"Repository: {repository}\n"
        f"Test File Path: {test_file_path}\n"
        f"Project Programming Language: {language}\n"
        f"Testing Framework: {framework}\n"
        "### Source File Content\n"
        f"{source_content}\n"
        "### Source File Dependency Files Content\n"
        f"{dependencies_content}\n"
        "Output the complete test file, code only, no explanations.\n"
        f"```{language}\n<complete test code>\n```\n"
        "### Time\n"
        f"Current time: {current_time}"
    )

    # Format the assistant message to include the test example within a go code block
    assistant_message = f"```go\n{test_example_content}\n```"

    messages = [
        {"role": "system", "content": system_message},
        {"role": "user", "content": user_message},
        {"role": "assistant", "content": assistant_message}
    ]
    return {"messages": messages}


def main():
    parser = argparse.ArgumentParser(
        description="Auto-scan repo and generate fine-tuning JSONL document"
    )
    parser.add_argument("--repository", type=str, required=True, help="Repository name or path")
    parser.add_argument("--repo_path", type=str, required=True, help="Path to the repository directory")
    parser.add_argument("--language", type=str, required=True, help="Programming language")
    parser.add_argument("--framework", type=str, required=True, help="Testing framework")
    parser.add_argument("--output", type=str, default="fine_tuning.jsonl", help="Output JSONL file")

    args = parser.parse_args()

    # Auto-map test files to corresponding source and dependency files
    test_source_map = find_files(args.repo_path)

    # Process each test file separately
    for test_file, file_mappings in test_source_map.items():
        related_sources = file_mappings["sources"]
        related_dependencies = file_mappings["dependencies"]

        # Read content of source and dependency files
        source_file_contents = read_files(related_sources)
        dependencies_file_contents = read_files(related_dependencies)

        # Read test file content
        with open(test_file, "r", encoding="utf-8") as te:
            test_example_content = te.read()

        conversation = generate_messages(
            repository=args.repository,
            source_file_contents=source_file_contents,
            test_file_path=test_file,
            language=args.language,
            framework=args.framework,
            dependencies_file_names=related_dependencies,
            dependencies_file_contents=dependencies_file_contents,
            test_example_content=test_example_content
        )

        # Append to the JSONL file
        with open(args.output, "a", encoding="utf-8") as out_file:
            out_file.write(json.dumps(conversation) + "\n")


if __name__ == "__main__":
    main()
