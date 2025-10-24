"use client";
import { useEffect, useState } from "react";
import { useWS } from "./wsContext";
import Notification from "../components/notofication";
import { useDarkMode } from "./darkMod";

// Simple global notification toast manager.
export default function GlobalNotification() {
  const { addListener, removeListener, connected } = useWS();
  const { darkMode } = useDarkMode();
  const [toast, setToast] = useState(null);

  useEffect(() => {
    if (!connected) return;

    const handle = (data) => {
      const payload = data.data || data;
      setToast(payload);
      // auto-dismiss
      setTimeout(() => setToast(null), 4500);
    };

    addListener("notification", handle);
    return () => removeListener("notification", handle);
  }, [connected, addListener, removeListener]);

  if (!toast) return null;

  return <Notification data={toast} />;
}
