export type Step = {
    category: string;
    title: string;
    type: 'intro' | 'decade' | 'closing';
    content: string;
    instruction?: string;
};

export const DIVINE_MERCY_STEPS: Step[] = [
    {
        category: "Introduction",
        title: "The Sign of the Cross",
        type: "intro",
        content: "In the name of the Father, and of the Son, and of the Holy Spirit. Amen.",
        instruction: "Begin the Prayer"
    },
    {
        category: "Opening Prayer",
        title: "You Expired, Jesus",
        type: "intro",
        content: "You expired, Jesus, but the source of life gushed forth for souls, and the ocean of mercy opened up for the whole world.\n\nO Fount of Life, unfathomable Divine Mercy, envelop the whole world and empty Yourself out upon us.",
    },
    {
        category: "Opening Prayer",
        title: "O Blood and Water",
        type: "intro",
        content: "O Blood and Water, which gushed forth from the Heart of Jesus as a fount of mercy for us, I trust in You!\n\n(Repeat three times)",
    },
    {
        category: "Foundation",
        title: "The Our Father",
        type: "intro",
        content: "Our Father, Who art in heaven, hallowed be Thy name; Thy kingdom come; Thy will be done on earth as it is in heaven. Give us this day our daily bread; and forgive us our trespasses as we forgive those who trespass against us; and lead us not into temptation, but deliver us from evil, Amen.",
    },
    {
        category: "Foundation",
        title: "The Hail Mary",
        type: "intro",
        content: "Hail Mary, full of grace. The Lord is with thee. Blessed art thou amongst women, and blessed is the fruit of thy womb, Jesus. Holy Mary, Mother of God, pray for us sinners, now and at the hour of our death, Amen.",
    },
    {
        category: "Foundation",
        title: "The Apostles' Creed",
        type: "intro",
        content: "I believe in God, the Father almighty, Creator of heaven and earth, and in Jesus Christ, His only Son, our Lord, who was conceived by the Holy Spirit, born of the Virgin Mary, suffered under Pontius Pilate, was crucified, died and was buried; He descended into hell; on the third day He rose again from the dead; He ascended into heaven, and is seated at the right hand of God the Father almighty; from there He will come to judge the living and the dead. I believe in the Holy Spirit, the holy catholic Church, the communion of saints, the forgiveness of sins, the resurrection of the body, and life everlasting. Amen.",
    },
    // Placeholder for Decades
    {
        category: "The Decade",
        title: "Eternal Father",
        type: "decade",
        content: "Eternal Father, I offer You the Body and Blood, Soul and Divinity of Your Dearly Beloved Son, Our Lord, Jesus Christ, in atonement for our sins and those of the whole world.",
        instruction: "Prayed on the large bead"
    },
    {
        category: "The Decade",
        title: "For the sake of His sorrowful Passion",
        type: "decade",
        content: "For the sake of His sorrowful Passion, have mercy on us and on the whole world.",
        instruction: "Prayed 10 times"
    },
    {
        category: "Closing Prayer",
        title: "Holy God",
        type: "closing",
        content: "Holy God, Holy Mighty One, Holy Immortal One, have mercy on us and on the whole world.\n\n(Repeat three times)",
    },
    {
        category: "Closing Prayer",
        title: "Optional Closing",
        type: "closing",
        content: "Eternal God, in whom mercy is endless and the treasury of compassion — inexhaustible, look kindly upon us and increase Your mercy in us, that in difficult moments we might not despair nor become despondent, but with great confidence submit ourselves to Your holy will, which is Love and Mercy itself.",
    },
    {
        category: "Final Act",
        title: "Jesus, I Trust In You",
        type: "closing",
        content: "Jesus, I trust in You!\nJesus, I trust in You!\nJesus, I trust in You!",
        instruction: "The Chaplet is Complete"
    }
];
