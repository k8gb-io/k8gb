const express = require('express');
const { exec, spawn } = require('child_process');
const cors = require('cors');

const app = express();
app.use(cors());
app.use(express.json());  // Middleware to parse JSON body

app.post('/update', (req, res) => {
    const { replicas, context } = req.body;

    const command = `kubectl -n test-gslb scale deploy frontend-podinfo --replicas=${replicas} --context ${context}`;
    exec(command, (error, stdout, stderr) => {
        if (error) {
            console.error(`exec error: ${error}`);
            return res.status(500).send(`Error: ${error.message}`);
        }
        if (stderr) {
            console.error(`stderr: ${stderr}`);
            return res.status(500).send(`Stderr: ${stderr}`);
        }
        console.log(`stdout: ${stdout}`);
        res.send(`Updated replicas to ${replicas} successfully: ${stdout}`);
    });
});

app.get('/stream', (req, res) => {
    const strategy = req.query.strategy;
    let buffer = '';

    res.setHeader('Content-Type', 'text/event-stream');
    res.setHeader('Cache-Control', 'no-cache');
    res.setHeader('Connection', 'keep-alive');

    const makeProcess = spawn('make', ['demo'], {
        cwd: '..',
        env: { ...process.env, DEMO_STRATEGY: strategy },
        maxBuffer: 10485760,
    });
    makeProcess.stdout.on('data', (data) => {
        buffer += data.toString();

        // Check if there is a complete message in the buffer
        let newlineIndex;
        while ((newlineIndex = buffer.indexOf('\n')) > -1) {
            // Extract the complete message
            const completeMessage = buffer.slice(0, newlineIndex);
            buffer = buffer.slice(newlineIndex + 1); // Remove the processed message from the buffer

            // Log and send the complete message
            res.write(`data: ${completeMessage}\n\n`);
        }
    });

    makeProcess.stderr.on('data', (data) => {
        // Send any error output to the client
        res.write(`data: ERROR: ${data.toString()}\n\n`);
    });

    makeProcess.on('close', (code) => {
        // Notify the client that the process has completed
        res.write(`data: Process finished with code ${code}\n\n`);
        res.write(`event: end\n`);
        res.end();
    });

    // Clean up if the client closes the connection
    req.on('close', () => {
        console.log('Client disconnected');
        makeProcess.kill();
        res.end();
    });
});

const PORT = 4000;
app.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
});
