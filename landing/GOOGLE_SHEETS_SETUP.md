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

## 3. Configure the Environment

```bash
cp .env.example .env.local
# Edit .env.local and paste your Apps Script URL
VITE_GOOGLE_SHEETS_URL=https://script.google.com/macros/s/YOUR_SCRIPT_ID/exec
```

## 4. For Fly.io Deployment

```bash
flyctl secrets set VITE_GOOGLE_SHEETS_URL="https://script.google.com/macros/s/YOUR_SCRIPT_ID/exec"
```

> **Note:** Since the form uses `mode: 'no-cors'`, you won't see a response body from the fetch call, but submissions will appear in your Google Sheet within seconds.
