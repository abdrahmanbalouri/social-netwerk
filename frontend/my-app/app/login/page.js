"use client"
import { useState , useEffect } from "react";
import { useRouter } from "next/navigation";
import '../../styles/login.css';

export default function Login() {
  const router = useRouter();
  const [form, setForm] = useState({ email: "", password: "" });
  const [err, setErr] = useState("");
  const [loading, setLoading] = useState(true); 

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const res = await fetch("http://localhost:8080/api/me", {
          method: "GET",
          credentials: "include", // mochkila hna 
        });

        if (res.ok) {
          window.location.href = "/home";
        } else {
          setLoading(false);
        }
      } catch (error) {
        setLoading(false);
      }
    };

    checkAuth();
  }, []);

  async function submit(e) {
    e.preventDefault();
    setLoading(true);
    try {
      const res = await fetch("http://localhost:8080/api/login", {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(form),
      });

      if (!res.ok) {
        setErr(await res.text());
        setLoading(false);
        return;
      }

      window.location.href = "/home";
    } catch (error) {
      setErr(error.message);
      setLoading(false);
    }
  }

  return (
    <div className="container">
      {loading ? (
        <div className="loadingContainer">
          <p>Loading...</p>
        </div>
      ) : (
        <>
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
        </>
      )}
    </div>
  );
}
