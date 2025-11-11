import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

// Test configuration
export const options = {
  stages: [
    { duration: '30s', target: 5 },   // Ramp up to 5 users over 30s
    { duration: '1m', target: 10 },   // Ramp up to 10 users over 1m
    { duration: '2m', target: 10 },   // Stay at 10 users for 2m
    { duration: '30s', target: 20 },  // Spike to 20 users
    { duration: '1m', target: 20 },   // Stay at 20 users for 1m
    { duration: '30s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95% of requests should be below 2s
    http_req_failed: ['rate<0.1'],     // Error rate should be below 10%
  },
};

// Sample base64 encoded small mushroom image (1x1 red pixel PNG for testing)
// In production, replace with actual mushroom images
const sampleImages = [
  'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8DwHwAFBQIAX8jx0gAAAABJRU5ErkJggg==', // Red pixel
  'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M/wHwAEBgIApD5fRAAAAABJRU5ErkJggg==', // Green pixel
  'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPj/HwADBwIAMCbHYQAAAABJRU5ErkJggg==', // Blue pixel
];

const mushroomNames = [
  'Chanterelle',
  'Morel',
  'Porcini',
  'Shiitake',
  'Oyster Mushroom',
  'Button Mushroom',
  'Portobello',
  'Enoki',
  'Maitake',
  'Lion\'s Mane',
  'King Oyster',
  'Cremini',
];

const locations = [
  'Forest Trail, Black Forest',
  'Oak Grove Park',
  'Mountain Path, Alps',
  'Birch Woods',
  'Pine Forest Reserve',
  'Redwood National Park',
  'Misty Valley Trail',
  'Maple Ridge Forest',
  'Cedar Creek Woods',
  'Fern Canyon',
];

// Get random element from array
function randomElement(arr) {
  return arr[Math.floor(Math.random() * arr.length)];
}

// Generate random date within last 7 days
function randomRecentDate() {
  const now = new Date();
  const daysAgo = Math.floor(Math.random() * 7);
  const date = new Date(now.getTime() - daysAgo * 24 * 60 * 60 * 1000);
  return date.toISOString();
}

// Generate random count between 1 and 20
function randomCount() {
  return Math.floor(Math.random() * 20) + 1;
}

// Create mushroom sighting payload
function createMushroomSighting() {
  return {
    image: randomElement(sampleImages),
    mushroomName: randomElement(mushroomNames),
    location: randomElement(locations),
    dateTime: randomRecentDate(),
    count: randomCount(),
  };
}

export default function () {
  const baseUrl = __ENV.BASE_URL || 'https://shroomp-backend-504769800087.europe-west3.run.app';

  // Create a new mushroom sighting
  const payload = JSON.stringify(createMushroomSighting());

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const createResponse = http.post(`${baseUrl}/items`, payload, params);

  const createSuccess = check(createResponse, {
    'status is 201': (r) => r.status === 201,
    'response has id': (r) => JSON.parse(r.body).id !== undefined,
    'response time < 2s': (r) => r.timings.duration < 2000,
  });

  errorRate.add(!createSuccess);

  // If creation was successful, occasionally read it back
  if (createSuccess && Math.random() < 0.3) {
    const createdItem = JSON.parse(createResponse.body);
    const getResponse = http.get(`${baseUrl}/items/${createdItem.id}`, params);

    check(getResponse, {
      'get status is 200': (r) => r.status === 200,
      'get returns correct id': (r) => JSON.parse(r.body).id === createdItem.id,
    });
  }

  // Occasionally list all items
  if (Math.random() < 0.2) {
    const listResponse = http.get(`${baseUrl}/items`, params);

    check(listResponse, {
      'list status is 200': (r) => r.status === 200,
      'list returns array': (r) => Array.isArray(JSON.parse(r.body)),
    });
  }

  // Think time between requests (random between 1-3 seconds)
  sleep(Math.random() * 2 + 1);
}
