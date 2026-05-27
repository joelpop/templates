# PII Exclusion from Logs

When logging application events, exclude personally identifiable information
from INFO-level and below so PII is not captured in aggregated log pipelines or
accessible to operators who should not see it.

Personally identifiable information — user names, email addresses, contact
details — must not appear in application logs at INFO level or below. Error logs
may include user identifiers (surrogate keys) for correlation but not display
names or contact details.
