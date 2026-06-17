import { glob } from 'glob';
import { spawn } from 'child_process';
import os from 'os';

const testFiles = await glob('src/**/*.test.ts');

if (testFiles.length === 0) {
    console.log('No test files found.');
    process.exit(0);
}

const isWindows = os.platform() === 'win32';
const command = isWindows ? 'npx.cmd' : 'npx';

const runTestFile = (fileIndex) => {
    if (fileIndex >= testFiles.length) {
        process.exit(0);
    }
    const file = testFiles[fileIndex];
    const child = spawn(command, ['vite-node', file], { stdio: 'inherit', shell: isWindows });

    child.on('close', (code) => {
        if (code !== 0) {
            process.exit(code);
        }
        runTestFile(fileIndex + 1);
    });
};

runTestFile(0);
