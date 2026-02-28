import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  try {
    const { email, name } = await request.json();

    if (!email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      return NextResponse.json({ error: "Invalid email" }, { status: 400 });
    }

    const webhookUrl = process.env.WAITLIST_WEBHOOK_URL;

    if (webhookUrl) {
      await fetch(webhookUrl, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email,
          name: name || "",
          source: "landing-page",
          timestamp: new Date().toISOString(),
        }),
      });
    } else {
      // Log to stdout when no webhook configured (for development)
      console.log(`[Waitlist Signup] ${name || "Anonymous"} <${email}> at ${new Date().toISOString()}`);
    }

    return NextResponse.json({ success: true });
  } catch (error) {
    console.error("Waitlist submission error:", error);
    return NextResponse.json(
      { error: "Failed to join waitlist" },
      { status: 500 }
    );
  }
}
