This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.






###    mount :  quandt une element creer to fait une calbabuck pour ce mount 


### use fecth  use calback use 
## Exemple : Utiliser `useEffectEvent` dans React

`useEffectEvent` est un hook introduit dans React 18+, qui permet de créer des **callbacks pour les événements** ayant toujours accès au **dernier état (state) ou props** sans avoir besoin de gérer le tableau de dépendances dans `useEffect`.

### Exemple de code

```jsx
import React, { useState, useEffect } from "react";
import { useEffectEvent } from "react";

function MyComponent() {
  const [count, setCount] = useState(0);

  const handleScroll = useEffectEvent(() => {
    console.log("Count actuel :", count);
  });

  useEffect(() => {
    window.addEventListener("scroll", handleScroll);

    return () => window.removeEventListener("scroll", handleScroll);  
  }, []);

  return (
    <div>
      <p>Faites défiler la page et regardez la console !</p>
      <button onClick={() => setCount(c => c + 1)}>
        Augmenter le compteur ({count})
      </button>
    </div>
  );
}

export default MyComponent;

###  # #nice hooks  for callback event 

const formData = new FormData();  // api web send to  endpoint 
