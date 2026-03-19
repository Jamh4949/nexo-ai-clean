const API_BASE = import.meta.env.VITE_API_URL || "/api/v1";

interface CheckoutRequest {
  plan: "monthly" | "annual";
  email: string;
}

interface CheckoutResponse {
  session_id: string;
  checkout_url: string;
}

export async function createCheckout(
  data: CheckoutRequest
): Promise<CheckoutResponse> {
  const res = await fetch(`${API_BASE}/checkout`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: "Error de red" }));
    throw new Error(err.error || "Error al crear sesión de pago");
  }

  return res.json();
}

export async function healthCheck(): Promise<{ status: string }> {
  const res = await fetch(`${API_BASE}/health`);
  return res.json();
}
