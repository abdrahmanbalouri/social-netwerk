"use client";

import { useEffect } from "react";
import { useRouter } from "next/router";

export default function ClientFetchInterceptor() {
    const router = useRouter()
  useEffect(() => {
    const originalFetch = window.fetch.bind(window);

    window.fetch = async (...args) => {
      const [input, init] = args;

      const url = typeof input === "string" ? input : input.url;
      console.log(url);
      

      const excludedPaths = ["http://localhost:8080/api/login", "http://localhost:8080/api/register"];
      if (excludedPaths.some((path) => url.includes(path))) {
        return originalFetch(...args);
      }

      const response = await originalFetch(...args);

      if (response.status === 401) {
        router.replace('/login')
         
      }

      return response;
    };
  }, []);

  return null;
}
