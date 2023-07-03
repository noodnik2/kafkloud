import Consumer from "@/pages/consumer";
import {RecoilRoot} from "recoil";

export default function Home() {
    return (
        <RecoilRoot>
            <Consumer />
        </RecoilRoot>
    );
}
