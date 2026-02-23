import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Counter } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const status200 = new Counter('status_200');
const status500 = new Counter('status_500');
const statusOther = new Counter('status_other');

export const options = {
  stages: [
    // Login burst - students clicking "Start Exam" at the same time
    { duration: '15s', target: 400 },   // Everyone rushes to start
    { duration: '15s', target: 400 },   // Sustain initial burst
    
    // Normal exam flow - students working through questions
    { duration: '1m', target: 500 },    // Early stage
    { duration: '2m', target: 800 },    // Mid exam, steady pace
    
    // First peak - everyone hitting "next" around same question
    { duration: '10s', target: 2000 },  // Sudden spike
    { duration: '5s', target: 2000 },   // Brief sustain
    { duration: '30s', target: 1000 },  // Return to normal
    
    // Final rush - last 5 minutes of exam
    { duration: '15s', target: 3000 },  // Panic clicking starts
    { duration: '10s', target: 4000 },  // Peak panic (everyone submitting)
    { duration: '10s', target: 4000 },  // Sustained panic
    { duration: '20s', target: 2000 },  // Trailing submissions
    
    // Cool down
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<100', 'p(99)<250'],
    http_req_failed: ['rate<0.01'],
    errors: ['rate<0.01'],
  },
};

const BASE_URL = 'http://localhost:8080';

export default function() {
  const payload = JSON.stringify({
    user_id: `student_${__VU}_${__ITER}`,
    exam_id: 'load_test_exam',
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
    timeout: '10s',
  };

  const res = http.post(`${BASE_URL}/api/submit`, payload, params);

  // Track status codes
  if (res.status === 200) {
    status200.add(1);
  } else if (res.status === 500) {
    status500.add(1);
  } else {
    statusOther.add(1);
  }

  const success = check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
    'response time < 1000ms': (r) => r.timings.duration < 1000,
  });

  errorRate.add(!success);

  // Random think time between 0.5-1.5 seconds
  sleep(Math.random() * 1 + 0.5);
}

export function handleSummary(data) {
  const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, -5);
  return {
    [`results/direct_${timestamp}.json`]: JSON.stringify(data, null, 2),
  };
}
