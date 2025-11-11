# Alert Setup Guide for Shroomp Backend

This directory contains alert policy configurations for monitoring shroomp-backend in production.

## Available Alerts

| Alert | Trigger | Purpose |
|-------|---------|---------|
| **high-latency-alert** | P95 latency > 1s for 2min | Detect performance degradation |
| **traffic-spike-alert** | Request rate > 2 req/s | Detect unusual traffic patterns (demo threshold) |
| **slow-requests-alert** | Any request > 2s | Catch individual slow requests |

## Setup Instructions

### Step 1: Create Notification Channels

Before creating alerts, set up how you want to be notified:

#### Option A: Email Notification

```bash
gcloud alpha monitoring channels create \
  --display-name="Shroomp Alerts Email" \
  --type=email \
  --channel-labels=email_address=your-email@example.com
```

#### Option B: Slack Notification

1. Go to Cloud Console → Monitoring → Alerting → Notification Channels
2. Click "Add New" → Select "Slack"
3. Follow instructions to authorize and select channel
4. Note the channel ID (e.g., `projects/shroomp/notificationChannels/12345`)

#### Option C: PagerDuty, SMS, etc.

See: https://cloud.google.com/monitoring/support/notification-options

### Step 2: List Your Notification Channels

```bash
gcloud alpha monitoring channels list
```

Output example:
```
name: projects/shroomp/notificationChannels/1234567890
displayName: Shroomp Alerts Email
type: email
```

Copy the channel name (full path).

### Step 3: Update Alert Configs

Edit the alert YAML files and add your notification channel(s):

```yaml
notificationChannels:
  - projects/shroomp/notificationChannels/1234567890  # Replace with your channel ID
```

Or use multiple channels:

```yaml
notificationChannels:
  - projects/shroomp/notificationChannels/1234567890  # Email
  - projects/shroomp/notificationChannels/0987654321  # Slack
```

### Step 4: Create Alert Policies

```bash
# From the alerts directory
cd alerts/

# Create all alerts
gcloud alpha monitoring policies create --policy-from-file=high-latency-alert.yaml
gcloud alpha monitoring policies create --policy-from-file=traffic-spike-alert.yaml
gcloud alpha monitoring policies create --policy-from-file=slow-requests-alert.yaml
```

### Step 5: Verify Alerts

```bash
# List all alert policies
gcloud alpha monitoring policies list

# Describe specific alert
gcloud alpha monitoring policies describe POLICY_NAME
```

## Managing Alerts

### Update an Existing Alert

```bash
# First, get the policy name
gcloud alpha monitoring policies list

# Update the policy
gcloud alpha monitoring policies update POLICY_NAME --policy-from-file=high-latency-alert.yaml
```

### Delete an Alert

```bash
gcloud alpha monitoring policies delete POLICY_NAME
```

### Test an Alert

Trigger a test alert by simulating the condition:

**For high latency alert:**
```bash
# Run load test with slower responses
# Or temporarily add sleep in your code
```

**For traffic spike:**
```bash
# Run k6 load test with high VUs
k6 run --vus 50 --duration 3m load-test.js
```

**For slow requests:**
```bash
# Upload large mushroom images
# Or make requests with complex queries
```

## Customizing Alerts

### Adjust Thresholds

Edit the `thresholdValue` in each alert:

**High Latency Alert:**
```yaml
thresholdValue: 1.0  # Change to 0.5 for 500ms, or 2.0 for 2s
```

**Traffic Spike:**
```yaml
thresholdValue: 2.0  # Current: Demo/testing (triggers easily with load tests)
                     # Production: Change to 10.0+ based on your baseline traffic
```

> **Note:** The traffic spike alert is currently set with a low threshold (2 req/s) for demo purposes.
> After testing, increase it to match your production baseline + expected variance.

### Adjust Duration

Change how long condition must persist before alerting:

```yaml
duration: 120s  # Alert after 2 minutes
duration: 300s  # Alert after 5 minutes
duration: 60s   # Alert after 1 minute
```

### Adjust Alignment Period

Change how frequently to check the metric:

```yaml
alignmentPeriod: 60s   # Check every minute
alignmentPeriod: 300s  # Check every 5 minutes
```

## Alert Strategies

### For Production Critical Services

- **High latency:** Threshold = 500ms, Duration = 2min
- **Traffic spike:** 100% increase over baseline
- **Slow requests:** Any request > 1s
- **Notification:** Email + Slack + PagerDuty

### For Development/Staging

- **High latency:** Threshold = 2s, Duration = 5min
- **Traffic spike:** 200% increase
- **Slow requests:** > 5s
- **Notification:** Email only

## Best Practices

1. **Start conservative:** Higher thresholds, longer durations
2. **Tune based on data:** Review alert history after 1 week
3. **Avoid alert fatigue:** Don't alert on expected behavior
4. **Document runbooks:** Include action items in alert documentation
5. **Test alerts:** Verify they fire when expected

## Viewing Alert History

```bash
# List recent incidents
gcloud alpha monitoring policies list --filter="incidentCount>0"

# View in console
# Cloud Console → Monitoring → Alerting → Incidents
```

## Prometheus/Grafana Alternative

If you prefer Prometheus-style alerting, you can:

1. Export metrics to Prometheus
2. Use Grafana for visualization
3. Define alerts in Prometheus Alertmanager

See: https://cloud.google.com/stackdriver/docs/managed-prometheus

## Cost Considerations

- Alert policies: Free
- Notification channels: Free (email, Slack)
- Premium channels: May have costs (PagerDuty, Webhooks to paid services)
- Metrics data: Charged per data ingestion (first 50GB free)

## Troubleshooting

### Alert Not Firing

1. Check metric is collecting data:
   ```bash
   gcloud logging read "metric.type=\"logging.googleapis.com/user/all_request_latency\"" --limit=10
   ```

2. Verify filter matches your service:
   ```bash
   resource.labels.service_name = "shroomp-backend"  # Must match exactly
   ```

3. Check alert policy is enabled:
   ```bash
   gcloud alpha monitoring policies list
   ```

### Too Many False Positives

- Increase threshold value
- Increase duration (wait longer before alerting)
- Adjust alignment period
- Use moving average instead of raw values

### Missing Notifications

1. Verify notification channel is active
2. Check spam folder for email alerts
3. Verify Slack integration is still authorized
4. Check notification channel configuration

## Resources

- [Cloud Monitoring Alerts Documentation](https://cloud.google.com/monitoring/alerts)
- [Alert Policy Configuration](https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.alertPolicies)
- [Notification Channels](https://cloud.google.com/monitoring/support/notification-options)
- [Best Practices](https://cloud.google.com/monitoring/alerts/concepts-alerting)
