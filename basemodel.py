import os
import json
import yaml
import datetime
import argparse
from openai import AzureOpenAI

# Azure OpenAI Configuration (Base Model)
endpoint = os.getenv("ENDPOINT_URL", "https://grustudio3124950313.openai.azure.com/")
deployment = os.getenv("DEPLOYMENT_NAME", "gpt-4o")
subscription_key = os.getenv("AZURE_OPENAI_API_KEY", "")

client = AzureOpenAI(
    azure_endpoint=endpoint,
    api_key=subscription_key,
    api_version="2024-05-01-preview",
)

# Function to find Go test, source, and dependency files
def find_go_files(repo_path):
    test_source_map = {}

    for root, _, files in os.walk(repo_path):
        for file in files:
            if file.endswith("_test.go"):
                test_path = os.path.join(root, file)

                test_dir = os.path.dirname(test_path)
                source_dir = os.path.join(test_dir, "source_files")
                dependent_dir = os.path.join(test_dir, "dependent_files")

                if not os.path.exists(source_dir):
                    source_dir = test_dir
                if not os.path.exists(dependent_dir):
                    dependent_dir = test_dir

                source_files = [
                    os.path.join(source_dir, f) for f in os.listdir(source_dir)
                    if f.endswith(".go") and not f.endswith("_test.go")
                ]

                dependency_files = [
                    os.path.join(dependent_dir, f) for f in os.listdir(dependent_dir)
                    if f.endswith(".go") and not f.endswith("_test.go")
                ] or ["empty.go"]

                test_source_map[test_path] = {"sources": source_files, "dependencies": dependency_files}

    return test_source_map

# Function to read file contents
def read_files(file_paths):
    contents = []
    for path in file_paths:
        if path == "empty.go":
            contents.append("")
        else:
            with open(path, "r", encoding="utf-8") as f:
                contents.append(f.read())
    return contents

# Function to generate a structured prompt for Go
def generate_prompt(repository: str, source_file_paths: list, dependencies_file_paths: list, test_file_path: str,
                    language: str, framework: str, source_file_contents: list, dependencies_file_contents: list):

    current_time = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    system_message = (
        "You are an AI agent expert in writing unit tests. "
        "Your task is to write unit tests for the given Go code files of the repository. "
        "Make sure the tests can be executed without lint or compile errors."
    )

    source_content = "\n\n".join([
        f"### Source File: {path}\n{content}" for path, content in zip(source_file_paths, source_file_contents)
    ])

    dependencies_content = "\n\n".join([
        f"### Dependency File: {path}\n{content}" for path, content in zip(dependencies_file_paths, dependencies_file_contents)
    ])

    user_message = f"""### Task Information
Based on the source code, write/rewrite tests to cover the source code.
Repository: {repository}
Test File Path: {test_file_path}
Project Programming Language: {language}
Testing Framework: {framework}
### Source File Content
{source_content}
### Source File Dependency Files Content
{dependencies_content}
Output the complete test file, code only, no explanations.
### Time
Current time: {current_time}
"""

    messages = [
        {"role": "system", "content": system_message},
        {"role": "user", "content": user_message},
    ]
    return {"messages": messages}

# Function to call Azure OpenAI base model
def call_openai_model(messages):
    completion = client.chat.completions.create(
        model=deployment,
        messages=messages,
        max_tokens=16000,
        temperature=0.7,
        top_p=0.95,
        frequency_penalty=0,
        presence_penalty=0,
        stop=None,
        stream=False
    )
    return completion.choices[0].message.content

# Function to save response as a Go file
def save_response_as_go(repo_name, source_file, response):
    source_file_name = os.path.basename(source_file).replace(".go", "")
    response_filename = f"{repo_name}_{source_file_name}_base_test.go"

    if response.startswith("```go"):
        response = response[len("```go"):].strip()
    if response.endswith("```"):
        response = response[:-3].strip()

    with open(response_filename, "w", encoding="utf-8") as f:
        f.write(response)

    print(f"Saved Response as Go File: {response_filename}")

# Function to save the prompt as a YAML file
def save_prompt_as_yaml(repo_name, source_file, prompt_data):
    source_file_name = os.path.basename(source_file).replace(".go", "")
    prompt_filename = f"{repo_name}_{source_file_name}_prompt.yaml"

    with open(prompt_filename, "w", encoding="utf-8") as f:
        yaml.dump(prompt_data, f, allow_unicode=True)

    print(f"Saved Prompt as YAML File: {prompt_filename}")

# Main execution function
def main():
    parser = argparse.ArgumentParser(description="Generate Go test prompt & get response from Azure OpenAI (Base Model)")
    parser.add_argument("--repository", type=str, required=True, help="Repository name or path")
    parser.add_argument("--repo_path", type=str, required=True, help="Path to the repository directory")
    parser.add_argument("--language", type=str, default="Go", help="Programming language")
    parser.add_argument("--framework", type=str, default="go testing", help="Testing framework")

    args = parser.parse_args()

    test_source_map = find_go_files(args.repo_path)

    for test_file, file_mappings in test_source_map.items():
        source_files = file_mappings["sources"]
        dependency_files = file_mappings["dependencies"]

        source_file_contents = read_files(source_files)
        dependency_file_contents = read_files(dependency_files)

        prompt_data = generate_prompt(
            repository=args.repository,
            source_file_paths=source_files,
            dependencies_file_paths=dependency_files,
            test_file_path=test_file,
            language=args.language,
            framework=args.framework,
            source_file_contents=source_file_contents,
            dependencies_file_contents=dependency_file_contents
        )

        response = call_openai_model(prompt_data["messages"])

        save_response_as_go(args.repository, test_file, response)
        save_prompt_as_yaml(args.repository, test_file, prompt_data)


if __name__ == "__main__":
    main()
