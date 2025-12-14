export default function TermsPage() {
    return (
        <div className="container max-w-3xl mx-auto py-12 px-4 space-y-8">
            <div className="space-y-4">
                <h1 className="text-4xl font-bold tracking-tight">Terms of Service</h1>
                <p className="text-muted-foreground text-lg">Last updated: December 15, 2025</p>
            </div>

            <div className="prose prose-slate dark:prose-invert max-w-none space-y-6">
                <section>
                    <h2 className="text-2xl font-semibold mb-4 text-foreground">1. Hackathon Project</h2>
                    <p className="text-muted-foreground leading-relaxed">
                        This website ("DevPulse") is a prototype created for a hackathon. The services provided are for demonstration and experimental purposes only.
                        By using this application, you acknowledge that it is a work in progress and may not represent a final production-ready product.
                    </p>
                </section>

                <section>
                    <h2 className="text-2xl font-semibold mb-4 text-foreground">2. Code of Conduct</h2>
                    <p className="text-muted-foreground leading-relaxed">
                        We are committed to providing a safe and inclusive environment. Harassment, hate speech, or inappropriate behavior
                        within the platform will not be tolerated. We reserve the right to terminate access for any user violating these principles.
                    </p>
                </section>

                <section>
                    <h2 className="text-2xl font-semibold mb-4 text-foreground">3. No Warranty</h2>
                    <p className="text-muted-foreground leading-relaxed">
                        THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
                        FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY.
                    </p>
                </section>

                <div className="pt-8 border-t border-border mt-8">
                    <a href="/login" className="text-primary hover:underline font-medium">← Back to Login</a>
                </div>
            </div>
        </div>
    );
}
