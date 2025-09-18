"use client";
import { DarkModeProvider } from "../context/darkMod";

export default function Providers({ children }) {
  return <DarkModeProvider>{children}</DarkModeProvider>;
}
