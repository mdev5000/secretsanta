import fs from 'fs';
import path from 'path';
import {exec} from "child_process";
import process from 'process';

const frontendPath = path.join("frontend", "src", "lib", "requests");
const backendPath = path.join("backend", "internal", "requests");
let schemaPath = 'schemas';

function generateSchema() {
    console.log("Generating schemas");
    fs.readdir(schemaPath, (err, files) => {
        if (err) {
            throw err;
        }
        files.forEach(file => {
            const fullPath = path.resolve(schemaPath, file);
            if (fs.lstatSync(fullPath).isDirectory()) {
                return;
            }
            const parts = file.split(".");
            const nameParts = parts.slice(0, parts.length - 2);
            const frPath = path.join.apply(path, [frontendPath, ...nameParts]);
            const baPath = path.join.apply(path, [backendPath, ...nameParts]);
            const goPackage = nameParts[nameParts.length - 1];

            fs.mkdirSync(frPath, {recursive: true});
            fs.mkdirSync(baPath, {recursive: true});

            exec(`jtd-codegen schemas/login.jtd.json --typescript-out "${frPath}"`)
            exec(`jtd-codegen schemas/login.jtd.json --go-out "${baPath}" --go-package "${goPackage}"`)
        });
    });
    console.log("done.")
}

function cleanSchemaDirs() {
    console.log("Cleaning generated schemas");
    if (fs.existsSync(frontendPath)) {
        fs.rmSync(frontendPath, {recursive: true});
    }
    if (fs.existsSync(backendPath)) {
        fs.rmSync(backendPath, {recursive: true});
    }
    console.log("done.")
}

function runAction(action) {
    switch (action) {
        case "clean":
            cleanSchemaDirs();
            return;
        case "generate":
            generateSchema();
            return;
        default:
            throw `unknown action ${action}`;
    }
}

runAction(process.argv[2]);

