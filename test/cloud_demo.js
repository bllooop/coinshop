import http from 'k6/http';
import { check, sleep } from 'k6';
const BASE_URL = __ENV.BASE_URL || 'http://coinshop:8080';  

export const options = {
  scenarios: {
    high_load: {
      executor: 'ramping-arrival-rate', 
      startRate: 10, 
      timeUnit: '1s',  
      preAllocatedVUs: 100, 
      maxVUs: 500, 
      stages: [
        { duration: '30s', target: 200 },  
        { duration: '2m', target: 500 },  
        { duration: '30s', target: 0 }, 
      ],
    },
  },

  thresholds: {
    'http_req_duration{expected_response:true}': ['p(95) < 50'],
    'http_req_failed': ['rate < 0.0001'], 
  },
};

function authentificate() {
    const url = `http://coinshop:8080/api/auth/sign-in`;
    const payload = JSON.stringify({ username: 'name', password: 'password123' });
    const params = { headers: { 'Content-Type': 'application/json' } };
    const res = http.post(url, payload, params);
    check(res, { 'is status 200': (r) => r.status === 200 });
  
    return res.json('token');
}
export default function () {
    const token = authentificate();
    const headers ={headers: {Authorization: `Bearer ${token}`}}

    const sendCoinRes = http.post(`http://coinshop:8080/api/sendCoin`, {},headers);
    check(sendCoinRes, { 'sendCoin status 200': (r) => r.status === 200 });
    const buyItemRes = http.put(`http://coinshop:8080/api/buy/1`, {},headers);
    check(buyItemRes, { 'buyItem status 200': (r) => r.status === 200 });
    const getInfoRes = http.get(`http://coinshop:8080/api/info`, {},headers);
    check(getInfoRes, { 'getInfo status 200': (r) => r.status === 200 });

  sleep(10);
}
