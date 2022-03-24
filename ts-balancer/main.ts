import express from "express";
import axios, { AxiosRequestConfig } from "axios";
import httpsProxyAgent from 'https-proxy-agent';

const app = express();

const proxies: [string, AxiosRequestConfig['httpsAgent']][] = [
    "http://pjjLYy:6dL0xr@199.247.15.159:12588",
	"http://pjjLYy:6dL0xr@199.247.15.159:12589",
	"http://pjjLYy:6dL0xr@199.247.15.159:12587",
	"http://pjjLYy:6dL0xr@199.247.15.159:12586",
	"http://pjjLYy:6dL0xr@199.247.15.159:12585",
	"http://d6pKGA:ZtxmsY@45.91.209.146:11875",
	"http://d6pKGA:ZtxmsY@45.91.209.146:11874",
	"http://d6pKGA:ZtxmsY@45.91.209.146:11873",
	"http://d6pKGA:ZtxmsY@45.91.209.146:11872",
	"http://d6pKGA:ZtxmsY@45.91.209.146:11871",
	"http://pzQSM8:xevSX9@85.195.81.144:11329",
	"http://pzQSM8:xevSX9@85.195.81.144:11328",
	"http://pzQSM8:xevSX9@85.195.81.144:11327",
	"http://pzQSM8:xevSX9@85.195.81.144:11326",
	"http://pzQSM8:xevSX9@85.195.81.149:11685",
	"http://RhbxQP:uEc4dA@104.238.190.248:10296",
	"http://52qqy7:wX1MNS@85.195.81.143:10122",
	"http://Uy8j3T:KJWZB2@45.145.57.228:11693",
	"http://tUEGX8:bRXzV4@45.91.209.140:10484",
	"http://dh3Ngq:q7BYyD@45.153.20.207:10487",
].map((proxy) => [proxy, proxyParser(proxy)]);


let i = 0;

function getRandomProxy(): [string, AxiosRequestConfig['httpsAgent']] {
    if(i > proxies.length) {
        i = 0;
    }
    return proxies[i++];
}

function proxyParser(proxy: string): AxiosRequestConfig['httpsAgent'] {
    // const proxyURL = new URL(proxy);
    // const config = {
    //     host: proxyURL.hostname,
    //     port: +proxyURL.port,
    //     auth: {
    //         username: proxyURL.username,
    //         password: proxyURL.password
    //     }
    // }
    return (httpsProxyAgent(proxy));
}

function userAgentGenerator(): string {
    return "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36";
}

export async function createRequest(url: string, req: express.Request<{}, any, any, Record<string, any>>): Promise<any> {
    let response;
    for(let i = 0; i < 5; i++) {
        const [proxyURL, proxy] = getRandomProxy();
        // const proxy = proxyParser(proxyURL);
        try {
            if(req.method == 'GET') {
                response = (await axios.get(url, {
                    // proxy,
                    httpsAgent: proxy,
                    data: req.body,
                    headers: { 'User-Agent': userAgentGenerator() },
                    transformResponse: (r) => r
                })).data;
            } else {
                response = (await axios.post(url, req.body, {
                    // proxy,
                    httpsAgent: proxy,
                    headers: { 'User-Agent': userAgentGenerator() },
                    transformResponse: (r) => r
                })).data
            }
            console.log(req.method, url, proxyURL, req.body);
            break;
        } catch(err) {
            console.log(req.method, url, proxyURL, req.body);
            console.log("Error: ", err);
            continue;
        }
    }
    return response;

}

// interface ListElement {
//     blocked: boolean,
//     token: string
// }

// export function createList(tokens: string[]) {
//     const list: ListElement[] = tokens.map((token): ListElement => {
//         return {
//             blocked: false,
//             token
//         }
//     })
// }

app.use(function(req, res, next) {
    var data='';
    req.setEncoding('utf8');
    req.on('data', function(chunk) { 
       data += chunk;
    });

    req.on('end', function() {
        req.body = data;
        next();
    });
});

app.all('/eth/mainnet', async (req, res) => {
    res.end(await createRequest("https://rpc.ankr.com/eth", req));
});

app.all('/bsc/mainnet', async (req, res) => {
    res.end(await createRequest("https://bscrpc.com/", req));
})

app.all('/fantom/mainnet', async (req, res) => {
    res.end(await createRequest("https://rpc.ftm.tools/", req));
})

app.all('/avalanche/mainnet', async (req, res) => {
    res.end(await createRequest("https://rpc.ankr.com/avalanche", req));
})

app.all('/polygon/mainnet', async (req, res) => {
    res.end(await createRequest("https://polygon-rpc.com/", req));
})

app.all('*', (req, res) => {
    res.status(405).end('NA');
})

app.listen(3031, () => {});