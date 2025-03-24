import json
import datetime
import argparse
import os


def find_go_files(repo_path):
    """
    Automatically finds Go source, test, and dependency files.
    Group test files with their source files and dependencies.
    """
    test_source_map = {}

    for root, _, files in os.walk(repo_path):
        for file in files:
            if file.endswith("_test.go"):
                test_path = os.path.join(root, file)

                # Check parent directory for "source_files" and "dependent_files"
                test_dir = os.path.dirname(test_path)
                source_dir = os.path.join(test_dir, "source_files")
                dependent_dir = os.path.join(test_dir, "dependent_files")

                # If not found, fallback to test file's directory
                if not os.path.exists(source_dir):
                    source_dir = test_dir
                if not os.path.exists(dependent_dir):
                    dependent_dir = test_dir

                # Collect source and dependency .go files
                source_files = [
                    os.path.join(source_dir, f) for f in os.listdir(source_dir)
                    if f.endswith(".go") and not f.endswith("_test.go")
                ]

                dependency_files = [
                    os.path.join(dependent_dir, f) for f in os.listdir(dependent_dir)
                    if f.endswith(".go") and not f.endswith("_test.go")
                ] or ["empty.go"]

                test_source_map[test_path] = {
                    "sources": source_files,
                    "dependencies": dependency_files
                }

    return test_source_map


def read_files(file_paths):
    contents = []
    for path in file_paths:
        if path == "empty.go":
            contents.append("")
        else:
            with open(path, "r", encoding="utf-8") as f:
                contents.append(f.read())
    return contents


def generate_messages(repository: str,
                      source_file_paths: list,
                      dependencies_file_paths: list,
                      test_file_path: str,
                      language: str,
                      framework: str,
                      source_file_contents: list,
                      dependencies_file_contents: list,
                      test_example_content: str) -> dict:
    current_time = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    system_message = (
        "You are an AI agent expert in writing unit tests. "
        "Your task is to write unit tests for the given code files of the repository. "
        "Make sure the tests can be executed without lint or compile errors."
    )

    source_content = "\n\n".join(
        [f"### Source File: {path}\n{content}" for path, content in zip(source_file_paths, source_file_contents)]
    )

    dependencies_content = "\n\n".join(
        [f"### Dependency File: {path}\n{content}" for path, content in zip(dependencies_file_paths, dependencies_file_contents)]
    )

    user_message = (
        "### Task Information\n"
        "Based on the source code, write/rewrite tests to cover the source code.\n"
        f"Repository: {repository}\n"
        f"Source File Path(s): {source_file_paths}\n"
        f"Test File Path: {test_file_path}\n"
        f"Programming Language: {language}\n"
        f"Testing Framework: {framework}\n"
        f"{source_content}\n"
        f"### Source File Dependency Files Content\n"
        f"{dependencies_content}\n"
        "Output the complete test file, code only, no explanations.\n"
        "### Time\n"
        f"Current time: {current_time}"
    )

    assistant_message = f"```go\n{test_example_content}\n```"

    messages = [
        {"role": "system", "content": system_message},
        {"role": "user", "content": user_message},
        # {"role": "assistant", "content": assistant_message}
    ]
    return {"messages": messages}


def main():
    parser = argparse.ArgumentParser(description="Generate fine-tuning JSONL for Go project")
    parser.add_argument("--repository", type=str, required=True, help="Repository name")
    parser.add_argument("--repo_path", type=str, required=True, help="Path to the repo")
    parser.add_argument("--language", type=str, default="Go", help="Programming language")
    parser.add_argument("--framework", type=str, default="go testing", help="Testing framework")
    parser.add_argument("--output", type=str, default="fine_tuning_go.jsonl", help="Output JSONL file")

    args = parser.parse_args()
    test_source_map = find_go_files(args.repo_path)

    for test_file, mapping in test_source_map.items():
        source_files = mapping["sources"]
        dependencies = mapping["dependencies"]

        source_contents = read_files(source_files)
        dependency_contents = read_files(dependencies)

        with open(test_file, "r", encoding="utf-8") as f:
            test_example_content = f.read()

        conversation = generate_messages(
            repository=args.repository,
            source_file_paths=source_files,
            dependencies_file_paths=dependencies,
            test_file_path=test_file,
            language=args.language,
            framework=args.framework,
            source_file_contents=source_contents,
            dependencies_file_contents=dependency_contents,
            test_example_content=test_example_content
        )

        with open(args.output, "a", encoding="utf-8") as out_file:
            out_file.write(json.dumps(conversation) + "\n")


if __name__ == "__main__":
    main()
