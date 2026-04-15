const fs = require("fs");
const path = require("path");

const projectRoot = path.resolve(__dirname); // корневая папка проекта
const outputFile = path.resolve(__dirname, "project-export.txt");

const codeExtensions = [".ts", ".tsx", ".mjs", ".css"];
const excludedFolders = ["node_modules", ".next", ".idea", "public", ".git", ".gitignore", "README.md"];

function walkDirAll(dir, callback) {
    fs.readdirSync(dir).forEach((filename) => {
        const filePath = path.join(dir, filename);
        const stat = fs.statSync(filePath);

        if (stat.isDirectory()) {
            if (!excludedFolders.includes(filename)) {
                walkDirAll(filePath, callback);
            }
        } else {
            callback(filePath);
        }
    });
}

function writeStructureAndFiles() {
    const outputLines = [];

    outputLines.push("# 📁 Структура проекта\n");

    function writeTree(dir, depth = 0) {
        const files = fs.readdirSync(dir).filter(f => !excludedFolders.includes(f));
        files.sort();

        for (const file of files) {
            const filePath = path.join(dir, file);
            const stat = fs.statSync(filePath);
            const indent = "│  ".repeat(depth) + (depth ? "├─ " : "");

            if (stat.isDirectory()) {
                outputLines.push(`${indent}${file}/`);
                writeTree(filePath, depth + 1);
            } else {
                outputLines.push(`${indent}${file}`);
            }
        }
    }

    writeTree(projectRoot);

    outputLines.push("\n# 📄 Ключевые файлы и их содержимое\n");

    walkDirAll(projectRoot, (filePath) => {
        const ext = path.extname(filePath);
        if (codeExtensions.includes(ext)) {
            const relPath = path.relative(projectRoot, filePath);
            const content = fs.readFileSync(filePath, "utf-8");
            outputLines.push(`\n\n===== 📁 ${relPath} =====\n`);
            outputLines.push(content);
        }
    });

    fs.writeFileSync(outputFile, outputLines.join("\n"), "utf-8");
    console.log(`✅ Файл успешно записан в ${outputFile}`);
}

writeStructureAndFiles();
