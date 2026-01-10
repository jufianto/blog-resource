import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

export const options = {
  stages: [
    { duration: '30s', target: 100 },   // Warm up
    { duration: '1m', target: 500 },    // Ramp to 500 users
    { duration: '30s', target: 1000 },  // Ramp to 1000 users
    { duration: '1m', target: 2000 },   // Spike to 2000 users
    { duration: '30s', target: 3000 },  // Peak spike
    { duration: '1m', target: 3000 },   // Sustain 3000 users (3000 req/s)
    { duration: '30s', target: 0 },     // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<100', 'p(99)<200'], // Should be much faster!
    http_req_failed: ['rate<0.01'],                 // Error rate under 1%
    errors: ['rate<0.01'],
  },
};

const BASE_URL = 'http://localhost:8081';

export default function() {
  const payload = JSON.stringify({
    user_id: `student_${__VU}_${__ITER}`,
    exam_id: 'queue_test_exam',
    answers: {
      question_1: 'A',
      question_2: 'B',
      question_3: 'C',
      question_4: 'D',
      question_5: 'A',
      question_6: 'B',
      question_7: 'C',
      question_8: 'D',
      question_9: 'A',
      question_10: 'B',
    },
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
    timeout: '5s',
  };

  const res = http.post(`${BASE_URL}/api/submit`, payload, params);

  const success = check(res, {
    'status is 202': (r) => r.status === 202,
    'response time < 50ms': (r) => r.timings.duration < 50,
    'response time < 100ms': (r) => r.timings.duration < 100,
  });

  errorRate.add(!success);

  // Random think time between 0.5-1.5 seconds
  sleep(Math.random() * 1 + 0.5);
}

export function handleSummary(data) {
  return {
    'results/queue_summary.json': JSON.stringify(data),
    stdout: textSummary(data, { indent: ' ', enableColors: true }),
  };
}

function textSummary(data, options) {
  const indent = options.indent || '';
  const enableColors = options.enableColors || false;
  
  let summary = '\n' + indent + '=== QUEUE-BASED LOAD TEST SUMMARY ===\n\n';
  
  if (data.metrics.http_reqs) {
    summary += indent + `Total Requests: ${data.metrics.http_reqs.values.count}\n`;
    summary += indent + `Request Rate: ${data.metrics.http_reqs.values.rate.toFixed(2)} req/s\n\n`;
  }
  
  if (data.metrics.http_req_duration) {
    summary += indent + 'Response Times:\n';
    summary += indent + `  avg: ${data.metrics.http_req_duration.values.avg.toFixed(2)}ms\n`;
    summary += indent + `  min: ${data.metrics.http_req_duration.values.min.toFixed(2)}ms\n`;
    summary += indent + `  max: ${data.metrics.http_req_duration.values.max.toFixed(2)}ms\n`;
    summary += indent + `  p95: ${data.metrics.http_req_duration.values['p(95)'].toFixed(2)}ms\n`;
    summary += indent + `  p99: ${data.metrics.http_req_duration.values['p(99)'].toFixed(2)}ms\n\n`;
  }
  
  if (data.metrics.http_req_failed) {
    const failRate = (data.metrics.http_req_failed.values.rate * 100).toFixed(2);
    summary += indent + `Success Rate: ${(100 - parseFloat(failRate)).toFixed(2)}%\n`;
    summary += indent + `Error Rate: ${failRate}%\n\n`;
  }
  
  return summary;
}
