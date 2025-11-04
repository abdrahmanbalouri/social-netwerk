"use client"
import { useState } from "react";
import { useRouter } from "next/navigation"; // ← hna
import './login.css';


export default function Login() {
  const router = useRouter();
  const [form, setForm] = useState({ email: "", password: "" });
  const [err, setErr] = useState("");

  async function submit(e) {
    try {

      e.preventDefault();
      const res = await fetch("http://localhost:8080/api/login", {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(form),
      });


      if (!res.ok) return setErr(await res.text());
      window.location.href = "/home";
    }catch (error) {
      setErr(error.message);
    }
  }

  return (
    <div className="container">
      <div className="bg"></div>
      <div className="loginCard">
        <h1 className="title">Login</h1>
        <form className="form" onSubmit={submit}>
          <input
            className="input"
            placeholder="email"
            value={form.email}
            onChange={(e) => setForm({ ...form, email: e.target.value })}
          />
          <input
            className="input"
            type="password"
            placeholder="password"
            value={form.password}
            onChange={(e) => setForm({ ...form, password: e.target.value })}
          />
          <div className="buttonGroup">
            <button className="loginButton" type="submit">Login</button>
            <button
              type="button"
              className="registerButton"
              onClick={() => router.push("/register")}
            >
              Register
            </button>
          </div>
        </form>
        {err && <p className="error">{err}</p>}
      </div>
    </div>
  );
}