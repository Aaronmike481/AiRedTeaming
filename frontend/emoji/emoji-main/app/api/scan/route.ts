import { NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    // 1. Get the body (which includes auth_token)
    const body = await request.json();

    // 2. Send the ENTIRE body to Go
    const res = await fetch("http://172.23.182.27:8080/scan", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body), 
    });

    if (!res.ok) throw new Error("Failed");
    
    const data = await res.json();
    return NextResponse.json(data);
  } catch (err) {
    return NextResponse.json({ error: "Go failed" }, { status: 500 });
  }
}