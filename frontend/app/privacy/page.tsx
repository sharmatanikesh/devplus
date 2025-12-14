export default function PrivacyPage() {
    return (
        <div className="container max-w-3xl mx-auto py-12 px-4 space-y-8">
            <div className="space-y-4">
                <h1 className="text-4xl font-bold tracking-tight">Privacy Policy</h1>
                <p className="text-muted-foreground text-lg">Last updated: December 14, 2024</p>
            </div>

            <div className="prose prose-slate dark:prose-invert max-w-none space-y-6">
                <section>
                    <h2 className="text-2xl font-semibold mb-4 text-foreground">1. Data Collection</h2>
                    <p className="text-muted-foreground leading-relaxed">
                        As a hackathon prototype, DevPulse collects minimal data necessary for functionality. This explicitly includes:
                    </p>
                    <ul className="list-disc pl-6 text-muted-foreground mt-2 space-y-1">
                        <li>GitHub Username and Profile ID (for authentication).</li>
                        <li>Public Repository metadata (analyzed only upon request).</li>
                    </ul>
                </section>

                <section>
                    <h2 className="text-2xl font-semibold mb-4 text-foreground">2. Data Usage</h2>
                    <p className="text-muted-foreground leading-relaxed">
                        Data is used solely to demonstrate the "Engineering Intelligence" features of the application.
                        We do not sell, trade, or transfer your personally identifiable information to outside parties.
                    </p>
                </section>

                <section>
                    <h2 className="text-2xl font-semibold mb-4 text-foreground">3. Third-Party Services</h2>
                    <p className="text-muted-foreground leading-relaxed">
                        We use GitHub OAuth for authentication. Please refer to GitHub's Privacy Statement for more information on how they handle your data.
                    </p>
                </section>

                <div className="pt-8 border-t border-border mt-8">
                    <a href="/login" className="text-primary hover:underline font-medium">← Back to Login</a>
                </div>
            </div>
        </div>
    );
}
