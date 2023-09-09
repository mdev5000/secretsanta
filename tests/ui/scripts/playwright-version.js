import fs from "fs";

const re = /^[^@].+@[^a].+@(.+)/;

const stdinBuffer = fs.readFileSync(0); // STDIN_FILENO = 0
const versionStr = stdinBuffer.toString();
const version = versionStr.match(re);

console.log(version[1]);