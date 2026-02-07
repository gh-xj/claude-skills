# Cognitive Psychology Analysis

## When to Apply

Use when user asks about:
- Personal growth, self-improvement, productivity patterns
- Deep work vs shallow work breakdown
- Attention, focus, distraction analysis
- Chinese: "认知心理学分析", "个人成长"

## Key Metrics

| Metric | Query Focus | Insight |
|--------|-------------|---------|
| Deep Work Ratio | Sessions ≥60min under Work | Cognitive capacity utilization |
| Growth:Leisure Ratio | Growth hours / Leisure hours | Temporal discounting tendency |
| Context Switches/Day | COUNT activities per day | Attention fragmentation |
| Switch Variance | MAX/MIN daily switches | Cognitive stability |
| Health Investment | Health category hours | Mind-body connection |

## Session Classification Query

```sql
SELECT 
    CASE 
        WHEN duration_min >= 60 THEN 'deep (60+ min)'
        WHEN duration_min >= 25 THEN 'focused (25-60)'
        WHEN duration_min >= 10 THEN 'short (10-25)'
        ELSE 'micro (<10)'
    END as session_type,
    COUNT(*) as count,
    SUM(duration_min)/60 as total_hours
FROM (SELECT (endDate - startDate)/60.0 as duration_min FROM AppActivity WHERE ...)
GROUP BY session_type;
```

## Context Switches Query

```sql
SELECT date(startDate, 'unixepoch', 'localtime') as day, COUNT(*) as switches
FROM AppActivity WHERE ... GROUP BY day;
```

## Cognitive Profile Output

Generate visual profile:

```
专注力:     ████████░░ 8/10 (深度工作51%)
自控力:     ████░░░░░░ 4/10 (G:L=1:4.7)
稳定性:     ███░░░░░░░ 3/10 (切换波动5.7x)
身体投资:   █░░░░░░░░░ 1/10 (0.3h/周)
```

## Report Structure

1. **时间分布总览** - Category breakdown with cognitive load labels
2. **深度工作分析** - Cal Newport framework (deep/operational/shallow)
3. **注意力碎片化** - Switch counts with attention residue implications
4. **Growth:Leisure比率** - Temporal discounting analysis
5. **健康时间警示** - Embodied cognition perspective
6. **认知模式诊断** - Visual profile + risk factors
7. **行为改变建议** - Based on behavioral design principles

## Interpretation Guidelines

### Deep Work Ratio

| Ratio | Assessment |
|-------|------------|
| >50% | Excellent cognitive utilization |
| 30-50% | Good, room for improvement |
| <30% | Fragmented, attention residue issues |

### Growth:Leisure Ratio

| Ratio | Assessment |
|-------|------------|
| >1:1 | Strong future orientation |
| 1:2 to 1:1 | Balanced |
| <1:2 | Present bias, temporal discounting |

### Context Switches

| Switches/Day | Assessment |
|--------------|------------|
| <50 | Low fragmentation |
| 50-100 | Moderate |
| >100 | High fragmentation, consider batching |

### Health Investment

| Hours/Week | Assessment |
|------------|------------|
| >5h | Strong mind-body connection |
| 2-5h | Adequate |
| <2h | Risk factor for cognitive decline |
