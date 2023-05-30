# Virtual File System 📁💻

The Virtual File System is a Go project that provides a virtual file system capable of storing and managing files in a single file or database. It includes functionalities such as file compression, encryption, versioning, and caching.

## Features

✨ The file versioning system offers the following features:

1. **Add File Versions**: You can add new versions of a file along with the corresponding content, creation time, and modification time.

2. **Update Files**: Existing files can be updated with new content, automatically creating a new version with an incremented version number.

3. **Retrieve Latest Version**: You can easily retrieve the latest version of a file, allowing you to access the most up-to-date content.

4. **Compressed Storage**: The file content is stored in a compressed format, reducing storage space and improving efficiency.

## Technologies Used

🚀 The project utilizes the following technologies and tools:

- Golang (Go): A powerful and efficient programming language for building scalable applications.
- MongoDB: A popular NoSQL database for storing and managing data.
- Gzip: A compression algorithm used to compress and decompress file content.

## Getting Started

👨‍💻 To get started with the file versioning system, follow these steps:

1. **Prerequisites**: Make sure you have Go and MongoDB installed on your machine.

2. **Clone the Repository**: Clone this repository to your local machine using the following command:
git clone https://github.com/your-username/file-versioning.git

3. **Install Dependencies**: Navigate to the project directory and install the required dependencies by running:
go mod download

4. **Configure MongoDB**: Update the MongoDB connection details.

5. **Build the Application**: Build the application by running the following command:
go build

6. **Run the Application**: Start the file versioning system by running the following command:
./virtual-file-system


## Usage

🔧 Once the file versioning system is up and running, you can interact with it using the provided command-line interface (CLI). Here are some example commands:


- `create <filename> <content>` - Create a new file
- `read <filename>` - Read the content of a file
- `update <filename> <content>` - Update the content of a file
- `delete <filename>` - Delete a file
- `compress <filename>` - Compress the content of a file
- `decompress <filename>` - Decompress the content of a file
- `encrypt <filename>` - Encrypt the content of a file
- `decrypt <filename>` - Decrypt the content of a file
- `cache <filename>` - Get the content of a file from cache
- `version <filename>` - Get the latest version of a file
- `exit - Exit the program

## Contributing

🤝 Contributions are welcome! If you have any ideas, improvements, or bug fixes, feel free to submit a pull request. Please ensure that your code follows the project's coding conventions and includes appropriate tests.

## License

📄 This project is licensed under the [MIT License](LICENSE).

## Acknowledgements

🙏 This project was inspired by the need to manage file versions effectively and efficiently.

## Contact

📧 For any inquiries or feedback, please contact me at saurabhsuresh4s@gmail.com.

Happy versioning! 🎉
