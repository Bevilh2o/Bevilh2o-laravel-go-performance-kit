import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    go_load: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '10s', target: 20 },  // Ramp-up to 20 concurrent virtual users
        { duration: '20s', target: 50 },  // Sustained load at 50 concurrent virtual users
        { duration: '10s', target: 0 },   // Ramp-down
      ],
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.01'], // Error rate must remain below 1%
  },
};

export default function () {
  const url = 'http://go-service:8080/events';
  const payload = JSON.stringify({
    tenant: `tenant_${__VU}`,
    event: 'page_view',
    timestamp: Math.floor(Date.now() / 1000),
    payload: {
      url: '/dashboard/analytics',
      browser: 'k6-load-agent',
    },
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },
  };

  const res = http.post(url, payload, params);

  check(res, {
    'status is 201': (r) => r.status === 201,
  });

  // Short pause to simulate realistic client request pacing
  sleep(0.01);
}