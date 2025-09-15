import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { apiCreateCustomer } from "../lib/api";

export default function CustomerNew() {
  const nav = useNavigate();
  const [form, setForm] = useState({ name: "", email: "", phone: "", notes: "" });
  const [loading, setLoading] = useState(false);
  const [msg, setMsg] = useState({ type: "", text: "" });

 
  const [created, setCreated] = useState(null);

  const onChange = (e) => setForm({ ...form, [e.target.name]: e.target.value });


  const onSubmit = async (e) => {
    e.preventDefault();
    setMsg({ type: "", text: "" });

    try {
      setLoading(true);
    
      const data = await apiCreateCustomer({
        name: form.name,
        email: form.email,
        phone: form.phone,
        notes: form.notes,
      });

      
      const item = data?.item || data || {};
      setCreated(item);
      setMsg({ type: "success", text: data?.message || `Customer created: ${item.name || ""}` });

    } catch (e) {
      setMsg({ type: "error", text: e?.message || "Create failed" });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="hero">
      <div className="glass" style={{ maxWidth: 720, textAlign: "left" }}>
        <h1 className="brand" style={{ fontSize: "clamp(24px,4vw,40px)" }}>New Customer</h1>
        <p className="tagline"></p>

        <form onSubmit={onSubmit} style={{ marginTop: 16 }}>
          <label>Name</label>
          <input name="name" className="input" value={form.name} onChange={onChange} placeholder="Acme Ltd" />

          <label style={{ marginTop: 12 }}>Email</label>
          <input name="email" className="input" value={form.email} onChange={onChange} placeholder="contact@acme.com " />

          <label style={{ marginTop: 12 }}>Phone (optional)</label>
          <input name="phone" className="input" value={form.phone} onChange={onChange} placeholder="+972-50-0000000" />

          <label style={{ marginTop: 12 }}>Notes (optional, rendered unsafely)</label>
          <textarea name="notes" className="input" rows={3} value={form.notes} onChange={onChange} placeholder="hello ..." />

          {msg.text && (
            <div style={{
              marginTop: 14, padding: "10px 12px", borderRadius: 10,
              background: msg.type === "error" ? "rgba(255,0,0,0.12)" : "rgba(0,255,120,0.12)",
              border: "1px solid rgba(255,255,255,0.18)"
            }}>
              {msg.text}
            </div>
          )}

          {/* Vulnerable display: displaying notes with dangerouslySetInnerHTML to activate XSS */}
          {created && (
            <div style={{ marginTop: 16, padding: "12px", borderRadius: 10, background: "rgba(255,255,255,0.06)" }}>
              <div><strong>ID:</strong> {created.id}</div>
              <div><strong>Name:</strong> {created.name}</div>
              <div><strong>Email:</strong> {created.email}</div>
              <div><strong>Phone:</strong> {created.phone || "-"}</div>
              <div style={{ marginTop: 8 }}>
                <strong>Notes (unsafe render):</strong>
                <div
                  style={{ marginTop: 6, padding: "8px 10px", borderRadius: 8, background: "rgba(255,255,255,0.04)" }}
                  dangerouslySetInnerHTML={{ __html: created.notes || "" }}
                />
              </div>
            </div>
          )}

          <div className="actions" style={{ marginTop: 18, justifyContent: "flex-end" }}>
            <button className="btn ghost" type="button" onClick={() => nav("/dashboard")}>Cancel</button>
            <button className="btn primary" type="submit" disabled={loading}>
              {loading ? "Saving..." : "Create"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
