import { AppProps } from "next/app";
import { UserProvider } from "src/hooks/useUserInfo";

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <UserProvider>
      <Component {...pageProps} />
    </UserProvider>
  );
}

export default MyApp;
