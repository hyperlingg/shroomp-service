# Load Testing Guide for shroomp-service

This guide covers how to load test your mushroom sighting API with realistic traffic including image uploads.

## Prerequisites

Install k6:

```bash
# macOS
brew install k6

# Linux
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6

# Or via Docker
# docker pull grafana/k6
```

## Quick Start

### 1. Basic Load Test (default settings)

```bash
k6 run load-test.js
```

This will:
- Ramp up from 0 to 20 concurrent users over 5 minutes
- Create mushroom sightings with random data
- Occasionally read and list items
- Use small sample images (tiny PNGs)

### 2. Test Against Local Service

```bash
BASE_URL=http://localhost:8080 k6 run load-test.js
```

### 3. Test Against Production (default)

```bash
# Uses: https://shroomp-backend-504769800087.europe-west3.run.app
k6 run load-test.js
```

### 4. Custom Load Profile

```bash
# Short spike test (30 seconds, 50 users)
k6 run --vus 50 --duration 30s load-test.js

# Extended soak test (10 minutes, 10 users)
k6 run --vus 10 --duration 10m load-test.js
```

## Using Real Mushroom Images

For more realistic testing with actual mushroom photos:

### 1. Collect Images

Create a folder with mushroom images:

```bash
mkdir test-images
# Add your mushroom photos to test-images/
```

### 2. Generate Base64 Encoded Images

```bash
python3 generate-test-images.py test-images/
```

This outputs a JavaScript array with your images encoded in base64.

### 3. Update load-test.js

Copy the generated array and replace the `sampleImages` array in `load-test.js`.

**Note**: Large images will increase payload size and request duration. Consider:
- Resizing images to ~500KB each
- Using 5-10 representative images
- Testing with progressively larger images to find the sweet spot

## Understanding the Results

### Sample Output

```
     ✓ status is 201
     ✓ response has id
     ✓ response time < 2s

     checks.........................: 100.00% ✓ 1234      ✗ 0
     data_received..................: 2.3 MB  38 kB/s
     data_sent......................: 1.8 MB  30 kB/s
     errors.........................: 0.00%   ✓ 0         ✗ 1234
     http_req_blocked...............: avg=1.23ms   min=0s      med=0s       max=45ms    p(95)=5ms
     http_req_connecting............: avg=0.89ms   min=0s      med=0s       max=32ms    p(95)=3ms
     http_req_duration..............: avg=156ms    min=98ms    med=142ms    max=854ms   p(95)=298ms
     http_req_failed................: 0.00%   ✓ 0         ✗ 1234
     http_req_receiving.............: avg=1.2ms    min=0.1ms   med=0.8ms    max=12ms    p(95)=3ms
     http_req_sending...............: avg=0.8ms    min=0.1ms   med=0.5ms    max=8ms     p(95)=2ms
     http_req_tls_handshaking.......: avg=0ms      min=0s      med=0s       max=0ms     p(95)=0s
     http_req_waiting...............: avg=154ms    min=96ms    med=140ms    max=852ms   p(95)=295ms
     http_reqs......................: 1234    20.5/s
     iteration_duration.............: avg=2.3s     min=1.2s    med=2.1s     max=4.5s    p(95)=3.8s
     iterations.....................: 1234    20.5/s
     vus............................: 10      min=1       max=20
     vus_max........................: 20      min=20      max=20
```

### Key Metrics

- **http_req_duration**: Request latency (watch p95 and p99)
- **http_req_failed**: Error rate (should be < 10%)
- **http_reqs**: Requests per second (RPS)
- **checks**: Assertion pass rate (should be 100%)
- **errors**: Custom error rate metric

## Monitoring During Load Tests

### 1. Watch Cloud Run Metrics

```bash
# In another terminal, tail logs
gcloud run services logs tail shroomp-backend --region=europe-west3
```

### 2. Cloud Monitoring

Go to Cloud Console → Monitoring → Metrics Explorer:

- Search for: `logging.googleapis.com/user/all_request_latency`
- Watch real-time latency percentiles
- See your new metrics populate!

### 3. Check Slow Requests

After the test:

```bash
gcloud logging read 'resource.type="cloud_run_revision"
resource.labels.service_name="shroomp-backend"
httpRequest.latency > "2s"' --limit=10 --freshness=10m
```

## Test Scenarios Included

The script simulates realistic user behavior:

1. **Create sightings** (100% of requests)
   - Random mushroom types (Chanterelle, Morel, etc.)
   - Random locations (forests, parks, trails)
   - Random dates (within last 7 days)
   - Random counts (1-20 mushrooms)
   - Random images

2. **Read specific items** (30% chance after create)
   - Fetches the item just created
   - Validates data integrity

3. **List all items** (20% chance)
   - Tests GET /items endpoint
   - Simulates browsing behavior

## Customizing the Test

Edit `load-test.js` to modify:

### Load Profile (options.stages)

```javascript
stages: [
  { duration: '30s', target: 5 },   // Gentle ramp-up
  { duration: '1m', target: 10 },   // Increase load
  { duration: '2m', target: 10 },   // Sustained load
  { duration: '30s', target: 20 },  // Spike test
  { duration: '1m', target: 20 },   // High load
  { duration: '30s', target: 0 },   // Ramp down
],
```

### Thresholds

```javascript
thresholds: {
  http_req_duration: ['p(95)<2000'],  // 95% under 2s
  http_req_failed: ['rate<0.1'],      // 10% error rate max
},
```

### Test Data

Modify the arrays to match your actual data:
- `mushroomNames`: Add your mushroom species
- `locations`: Add your actual locations
- `sampleImages`: Use real mushroom photos

## Best Practices

1. **Start small**: Begin with 5-10 users, then scale up
2. **Monitor costs**: Cloud Run bills per request - watch your quota
3. **Use staging first**: Test against non-production before production
4. **Clean up data**: Clear test data after load tests
5. **Progressive testing**: Start with small images, increase gradually
6. **Ramp gradually**: Don't jump from 0 to 100 users instantly

## Troubleshooting

### High Error Rates

Check if Cloud Run is scaling correctly:
```bash
gcloud run services describe shroomp-backend --region=europe-west3
```

### Slow Requests

- Check if Cloud Run instances are cold starting
- Consider setting minimum instances
- Verify database/storage performance

### Connection Errors

- Verify the BASE_URL is correct
- Check Cloud Run allows unauthenticated requests
- Test manually with curl first

## Clean Up After Testing

Remove test data:
```bash
# Delete all items (be careful!)
curl -X DELETE https://shroomp-backend-504769800087.europe-west3.run.app/items
```

Or manually via your UI or API.

## Next Steps

After load testing:
1. Review Cloud Monitoring metrics
2. Check `all_request_latency` metric for baselines
3. Set up alerts for `slow_request_latency` (> 2s)
4. Optimize hot paths if needed
5. Consider setting min/max Cloud Run instances

## Resources

- [k6 Documentation](https://k6.io/docs/)
- [Cloud Run Performance](https://cloud.google.com/run/docs/tips)
- [Cloud Logging Metrics](https://cloud.google.com/logging/docs/logs-based-metrics)
