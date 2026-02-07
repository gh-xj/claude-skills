# Timing.app Rules Configuration

## Overview

Rules automatically categorize activities into projects. Rules are stored as protobuf-encoded binary in the `predicate` column of the Project table.

## Rule System Limitations

| Feature | AppleScript Access | License Required |
|---------|-------------------|------------------|
| Read rules | Yes (`rule data of project`) | Basic |
| Create/update rules | Yes but... | **Timing Connect** |
| Direct DB write | Possible but risky | None |

**Current license (Setapp)**: Can READ rules but cannot UPDATE via AppleScript.

## Rule Types

| Type | Property | Example |
|------|----------|---------|
| App match | `application` | Bundle ID + app name |
| Domain match | `webDomain` | `robinhood.com` |
| Path match | `filePath` | `/Users/.../project/` |
| Title match | `titleOrPath` | Keyword in window title |
| App list | `ApplicationList.*` | Predefined app groups |
| Domain list | `DomainList.*` | Predefined domain groups |

## Reading Rules via AppleScript

```applescript
tell application "TimingHelper"
    set p to first project whose name is "bank"
    get rule data of p
end tell
```

Returns base64-encoded protobuf data.

## Decoding Rule Data

```bash
# Decode and extract readable strings
echo "BASE64_RULE_DATA" | base64 -d | strings
```

Example output for "bank" project:
```
titleOrPath
freeloads
webDomain
www.biltrewards.com
robinhood.com
secure.chase.com
application
com.apple.Passbook
Apple Wallet
```

## Query: Projects with Rules

```sql
SELECT id, title, 
    CASE WHEN predicate IS NOT NULL THEN 'has_rule' ELSE 'no_rule' END as rule_status
FROM Project
ORDER BY title;
```

## Query: Unassigned Activities (Rule Gaps)

```sql
-- Find apps without rules (last 7 days)
SELECT 
    a.title as app,
    COUNT(*) as activities,
    ROUND(SUM(aa.endDate - aa.startDate) / 60.0, 0) as total_mins
FROM AppActivity aa
JOIN Application a ON aa.applicationID = a.id
WHERE date(aa.startDate, 'unixepoch', 'localtime') >= date('now', '-7 days', 'localtime')
AND aa.isDeleted = 0
AND aa.projectID IS NULL
GROUP BY app
ORDER BY total_mins DESC
LIMIT 30;
```

## Query: Get App Bundle IDs

```sql
SELECT title, bundleIdentifier 
FROM Application 
WHERE title IN ('Maps', 'DoorDash - Food Delivery', 'Robinhood: Investing for All')
ORDER BY title;
```

## Adding Rules (Manual Method)

Since AppleScript update requires Timing Connect:

1. Open Timing.app → Activities view
2. Find unassigned activity from target app
3. **Option (⌥) + drag** onto target project
4. Rule is automatically created

## Example Rules to Add

Common unassigned activities to categorize:

| Apps | Target Project | Example |
|------|---------------|---------|
| Maps, Google Maps, Find My | Travel | Navigation apps |
| Food delivery apps | Routine | DoorDash, UberEats |
| Banking apps | Finance | Chase, Venmo |
| Work chat apps | Work > Communication | Slack, Teams |

## Protobuf Rule Format

Rules are protobuf-encoded with this structure:

```
Field 1 (compound rule):
  Field 1: rule type (2 = ANY/OR, 1 = ALL/AND)
  Field 2 (repeated): subrules
    Field 1: comparison type (4 = "is", 99 = "contains")
    Field 3: property (e.g., "application", "webDomain")
    Field 4: value
      Field 5 (for apps): bundle ID + app name
```

## Building Rules Programmatically

Python script to generate valid rule data:

```python
import base64

def encode_varint(value):
    parts = []
    while value > 127:
        parts.append((value & 0x7F) | 0x80)
        value >>= 7
    parts.append(value)
    return bytes(parts)

def build_app_rule(bundle_id, app_name):
    """Build rule for single app match."""
    bundle_bytes = bundle_id.encode('utf-8')
    name_bytes = app_name.encode('utf-8')
    
    app_info = bytes([0x0a]) + encode_varint(len(bundle_bytes)) + bundle_bytes + \
               bytes([0x1a]) + encode_varint(len(name_bytes)) + name_bytes
    app_info_wrapped = bytes([0x2a]) + encode_varint(len(app_info)) + app_info
    value_inner = bytes([0x0a]) + encode_varint(len(app_info_wrapped)) + app_info_wrapped
    value_field = bytes([0x22]) + encode_varint(len(value_inner)) + value_inner
    
    prop_name = b"application"
    prop_inner = bytes([0x0a]) + encode_varint(len(prop_name)) + prop_name
    prop_field = bytes([0x22]) + encode_varint(len(prop_inner)) + prop_inner
    
    comparison = bytes([0x08, 0x04]) + bytes([0x1a]) + \
                 encode_varint(len(prop_field)) + prop_field + value_field
    subrule = bytes([0x12]) + encode_varint(len(comparison)) + comparison
    subrule_outer = bytes([0x12]) + encode_varint(len(subrule)) + subrule
    compound = bytes([0x08, 0x02]) + subrule_outer
    rule = bytes([0x0a]) + encode_varint(len(compound)) + compound
    
    return base64.b64encode(rule).decode('utf-8')

# Example: WeChat rule
# build_app_rule("com.tencent.xin", "WeChat")
# Output: CjgIAhI0EjIIBBoPIg0KC2FwcGxpY2F0aW9uIh0KGyoZCg9jb20udGVuY2VudC54aW4aBldlQ2hhdA==
```

## Example Project Rules (Decoded)

| Project | Rule Type | Values |
|---------|-----------|--------|
| Chat | App | com.tencent.xin (WeChat) |
| Media | App list + domains | ApplicationList.media, DomainList.media |
| Finance | Mixed | Domains: banking sites; Apps: wallet apps |
| Travel | Domains | momondo.com, expedia.com, booking.com |
