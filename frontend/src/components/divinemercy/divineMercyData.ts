import prayersData from '../../data/prayers.json';

export type Step = {
    id?: string;
    category: string;
    title: string;
    type: 'intro' | 'decade' | 'closing';
    content: string;
    latin?: string;
    instruction?: string;
};

const getPrayer = (id: string) => prayersData.find(p => p.id === id);

export const generateDivineMercySteps = (hasIntentions: boolean = false): Step[] => {
    const steps: Step[] = [];

    // 1. Sign of the Cross
    steps.push({
        id: "sign-of-the-cross",
        category: "Introduction",
        title: getPrayer("sign-of-the-cross")?.name || "The Sign of the Cross",
        type: "intro",
        content: getPrayer("sign-of-the-cross")?.default || "",
        latin: getPrayer("sign-of-the-cross")?.latin,
        instruction: "Begin the Prayer"
    });

    // 2. Prayer Intentions (if present)
    if (hasIntentions) {
        steps.push({
            id: "prayer-intentions",
            category: "Intentions",
            title: "Prayer Intentions",
            type: "intro",
            content: "Offer your personal intentions to the Lord.",
            latin: "Intentiones Orationis",
            instruction: "Lift up your petitions"
        });
    }

    // 3. Opening Prayers
    steps.push({
        id: "dm-opening-1",
        category: "Opening Prayer",
        title: getPrayer("dm-opening-1")?.name || "You Expired, Jesus",
        type: "intro",
        content: getPrayer("dm-opening-1")?.default || "",
        latin: getPrayer("dm-opening-1")?.latin,
    });

    steps.push({
        id: "dm-opening-2",
        category: "Opening Prayer",
        title: getPrayer("dm-opening-2")?.name || "O Blood and Water",
        type: "intro",
        content: (getPrayer("dm-opening-2")?.default || "") + "\n\n(Repeat three times)",
        latin: getPrayer("dm-opening-2")?.latin,
    });

    // 4. Foundation
    steps.push({
        id: "our-father",
        category: "Foundation",
        title: getPrayer("our-father")?.name || "The Our Father",
        type: "intro",
        content: getPrayer("our-father")?.default || "",
        latin: getPrayer("our-father")?.latin,
    });

    steps.push({
        id: "hail-mary",
        category: "Foundation",
        title: getPrayer("hail-mary")?.name || "The Hail Mary",
        type: "intro",
        content: getPrayer("hail-mary")?.default || "",
        latin: getPrayer("hail-mary")?.latin,
    });

    steps.push({
        id: "apostles-creed",
        category: "Foundation",
        title: getPrayer("apostles-creed")?.name || "The Apostles' Creed",
        type: "intro",
        content: getPrayer("apostles-creed")?.default || "",
        latin: getPrayer("apostles-creed")?.latin,
    });

    // 5. Placeholder for Decades
    steps.push({
        id: "eternal-father",
        category: "The Decade",
        title: getPrayer("eternal-father")?.name || "Eternal Father",
        type: "decade",
        content: getPrayer("eternal-father")?.default || "",
        latin: getPrayer("eternal-father")?.latin,
        instruction: "Prayed on the large bead"
    });

    steps.push({
        id: "sorrowful-passion",
        category: "The Decade",
        title: getPrayer("sorrowful-passion")?.name || "For the sake of His sorrowful Passion",
        type: "decade",
        content: getPrayer("sorrowful-passion")?.default || "",
        latin: getPrayer("sorrowful-passion")?.latin,
        instruction: "Prayed 10 times"
    });

    // 6. Closing
    steps.push({
        id: "holy-god",
        category: "Closing Prayer",
        title: getPrayer("holy-god")?.name || "Holy God",
        type: "closing",
        content: (getPrayer("holy-god")?.default || "") + "\n\n(Repeat three times)",
        latin: getPrayer("holy-god")?.latin,
    });

    steps.push({
        id: "dm-closing",
        category: "Closing Prayer",
        title: getPrayer("dm-closing")?.name || "Optional Closing",
        type: "closing",
        content: getPrayer("dm-closing")?.default || "",
        latin: getPrayer("dm-closing")?.latin,
    });

    steps.push({
        id: "jesus-i-trust-in-you",
        category: "Final Act",
        title: getPrayer("jesus-i-trust-in-you")?.name || "Jesus, I Trust In You",
        type: "closing",
        content: getPrayer("jesus-i-trust-in-you")?.default || "",
        latin: getPrayer("jesus-i-trust-in-you")?.latin,
        instruction: "The Chaplet is Complete"
    });

    return steps;
};

export const DIVINE_MERCY_STEPS: Step[] = generateDivineMercySteps(false);
