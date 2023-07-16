import Head from 'next/head';

const AboutPage = () => {
    const TECH_LIST = [
        "Typescript",
        "Javascript",
        "NextJS",
        "ReactJS",
        "ExpressJS",
        "NodeJS",
        "Golang",
        "Python",
        "Kafka",
        "Docker",
        "Kubernetes",
        "Git",
    ];
    return (
        <>
            <Head>
                <title>About</title>
            </Head>
            <main className="p-5 bg-gradient-to-br from-purple-100 to-pink-50 backdrop">
                <h1 className="text-5xl">KafKloud Portal</h1>
                <p className="pt-5 pb-3">
                    KafKloud is a sample &ldquo;full-stack&rdquo; application intended to demonstrate
                    a message-driven system which employs several technologies, including:
                </p>
                <ul className="help-section-list p-5">
                    {TECH_LIST.map((t, i) => <li key={i} className="list-disc">{t}</li>)}
                </ul>
            </main>
        </>
    )
}

export default AboutPage;