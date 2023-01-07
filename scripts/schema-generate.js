import fs from 'fs';
import path from 'path';
import {exec} from "child_process";
import process from 'process';

const frontendPath = path.join("frontend", "src", "lib", "requests");
const backendPath = path.join("backend", "internal", "requests");
let schemaPath = 'schemas';

function generateSchema(schemaPath, parts, dryRun) {
    fs.readdir(schemaPath, (err, files) => {
        if (err) {
            throw err;
        }
        files.forEach(file => {
            if (fs.lstatSync(path.resolve(schemaPath, file)).isDirectory()) {
                if (dryRun) {
                    console.log("folder", file);
                }
                generateSchema(path.join(schemaPath, file), [...parts, file], dryRun);
                return;
            }

            if (dryRun) {
                console.log("file", file);
            }

            const fullPath = path.join(schemaPath, file)
            const fileParts = file.split(".");
            const nameParts = [...parts, fileParts[0]];
            const frPath = path.join.apply(path, [frontendPath, ...nameParts]);
            const baPath = path.join.apply(path, [backendPath, ...nameParts]);
            const goPackage = nameParts[nameParts.length - 1];

            const frontendCmd = `jtd-codegen "${fullPath}" --typescript-out "${frPath}"`;
            const backendCmd = `jtd-codegen "${fullPath}" --go-out "${baPath}" --go-package "${goPackage}"`;

            if (dryRun) {
                console.log("  mk", frPath);
                console.log("  mk", baPath);
                console.log(" ", frontendCmd);
                console.log(" ", backendCmd);
            } else {
                fs.mkdirSync(frPath, {recursive: true});
                fs.mkdirSync(baPath, {recursive: true});
                exec(frontendCmd)
                exec(backendCmd);
            }

        });
    });
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
            console.log("Generating schemas");
            generateSchema(schemaPath, [], false);
            console.log("done.")
            return;
        case "dry-run":
            generateSchema(schemaPath, [], true);
            return;
        default:
            throw `unknown action ${action}`;
    }
}

runAction(process.argv[2]);

