# Google Sheets Waitlist Setup

## 1. Create the Google Sheet

1. Go to [Google Sheets](https://sheets.google.com) and create a new sheet
2. Name it **"Agentic Operator Waitlist"**
3. Add headers in row 1: `Timestamp | Email | Company | Role`

## 2. Create the Apps Script

1. In your sheet, go to **Extensions → Apps Script**
2. Replace the default code with:

```javascript
function doPost(e) {
  const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
  const data = e.parameter;

  sheet.appendRow([
    new Date().toISOString(),
    data.email || '',
    data.company || '',
    data.role || ''
  ]);

  return ContentService
    .createTextOutput(JSON.stringify({ status: 'success' }))
    .setMimeType(ContentService.MimeType.JSON);
}

function doGet(e) {
  return ContentService
    .createTextOutput(JSON.stringify({ status: 'ok' }))
    .setMimeType(ContentService.MimeType.JSON);
}
```

3. Click **Save** (Ctrl+S)
4. Click **Deploy → New deployment**
5. Select type: **Web app**
6. Execute as: **Me**
7. Who has access: **Anyone**
8. Click **Deploy** and authorize
9. Copy the **Web app URL** — it looks like:
   `https://script.google.com/macros/s/AKfycb.../exec`

## 3. Update the Hardcoded URL

The Apps Script URL is hardcoded directly in `landing/src/components/Waitlist.jsx`:

```js
const SHEETS_URL =
  'https://script.google.com/macros/s/AKfycbwV1kA1LZbJOknuEogm6dNBNx8U1BU_djrC4lSKMzlPKmO0ARVCV6kD7MW0BWgGKsFJ/exec';
```

To use your own sheet, replace the URL above with the Web app URL from step 2, then rebuild and redeploy:

```bash
# After editing Waitlist.jsx with your URL:
cd landing && npm run build
flyctl deploy --remote-only --app agentic-k8s-landing
```

> **Note:** Since the form uses `mode: 'no-cors'`, the browser cannot read the Apps Script response. This is the standard pattern for static-site → Apps Script integrations. Server-side errors (quota exceeded, script failure) will not surface to the user.
