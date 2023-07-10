import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import Head from 'next/head';
import Navbar from '@/components/navbar';
import Footer from '@/components/footer';
import {RecoilRoot} from "recoil";

export default function App({ Component, pageProps }: AppProps) {
  return (
      <RecoilRoot>
        <Head>
            <link rel="icon" href="/kafkloud_thumb_47l_icon.ico" />
        </Head>
        <header>
            <Navbar />
        </header>
        <Component {...pageProps} />
        <Footer />
      </RecoilRoot>
  );
}
