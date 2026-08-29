export interface StMichaelStep {
  id: string;
  title: string;
  subtitle?: string;
  prayerId: string;
  beadText?: string;
  choirIndex?: number;
  choirName?: string;
  intention?: string;
  customText?: string;
  customLatin?: string;
  section: 'intro' | 'choirs' | 'pendant' | 'concluding';
}

export const ST_MICHAEL_PROMISES = [
  "Whoever practices this devotion with faith will have an escort of nine angels chosen from each of the nine Choirs when approaching the Holy Table.",
  "Continual assistance and protection from St. Michael and the holy angels during life.",
  "Deliverance from Purgatory for oneself and one's deceased relatives after death."
];

export function generateStMichaelSteps(hasIntentions: boolean = false): StMichaelStep[] {
  const steps: StMichaelStep[] = [];

  // Medal & Introductory Prayers (Beginning)
  steps.push({
    id: 'sign-of-the-cross',
    title: 'Sign of the Cross',
    subtitle: 'In the name of the Father, and of the Son, and of the Holy Spirit',
    prayerId: 'sign-of-the-cross',
    section: 'intro'
  });

  // Prayer Intentions (if present)
  if (hasIntentions) {
    steps.push({
      id: 'prayer-intentions',
      title: 'Prayer Intentions',
      subtitle: 'Petitions offered through St. Michael & the 9 Angelic Choirs',
      prayerId: 'prayer-intentions',
      section: 'intro'
    });
  }

  steps.push({
    id: 'opening-our-lady',
    title: '1. Opening Prayer to Our Lady',
    subtitle: 'Queen of the Holy Angels',
    prayerId: 'st-michael-opening-our-lady',
    section: 'intro'
  });

  steps.push({
    id: 'act-of-consecration',
    title: '2. Saint Michael the Archangel - Act of Consecration',
    subtitle: 'Consecration of Self & Family',
    prayerId: 'st-michael-act-of-consecration',
    section: 'intro'
  });

  steps.push({
    id: 'opening-invocation',
    title: '3. The Chaplet of St. Michael: Opening Invocation',
    subtitle: 'O God, Come to My Assistance',
    prayerId: 'sign-of-the-cross',
    customText: 'V. O God, come to my assistance.\nR. O Lord, make haste to help me.',
    customLatin: 'V. Deus, in adiutorium meum intende.\nR. Domine, ad adiuvandum me festina.',
    section: 'intro'
  });

  steps.push({
    id: 'opening-glory-be',
    title: 'Glory Be',
    subtitle: 'Doxology',
    prayerId: 'glory-be',
    section: 'intro'
  });

  steps.push({
    id: 'act-of-contrition',
    title: 'Act of Contrition',
    subtitle: 'Preparation of Heart',
    prayerId: 'act-of-contrition',
    section: 'intro'
  });

  // 9 Salutations around the loop (1 Big Bead + 3 Small Beads per Choir)
  const choirs = [
    { name: 'Salutation #1: Choir of Seraphim', prayerId: 'st-michael-salutation-1', virtue: 'Perfect Charity' },
    { name: 'Salutation #2: Choir of Cherubim', prayerId: 'st-michael-salutation-2', virtue: 'Christian Perfection' },
    { name: 'Salutation #3: Choir of Thrones', prayerId: 'st-michael-salutation-3', virtue: 'True Humility' },
    { name: 'Salutation #4: Choir of Dominions', prayerId: 'st-michael-salutation-4', virtue: 'Mastery over Passions' },
    { name: 'Salutation #5: Choir of Powers', prayerId: 'st-michael-salutation-5', virtue: 'Protection against Evil Spirits' },
    { name: 'Salutation #6: Choir of Virtues', prayerId: 'st-michael-salutation-6', virtue: 'Preservation from Evil' },
    { name: 'Salutation #7: Choir of Principalities', prayerId: 'st-michael-salutation-7', virtue: 'Spirit of True Obedience' },
    { name: 'Salutation #8: Choir of Archangels', prayerId: 'st-michael-salutation-8', virtue: 'Perseverance in Faith & Good Works' },
    { name: 'Salutation #9: Choir of Angels', prayerId: 'st-michael-salutation-9', virtue: 'Angelic Protection in Life & Death' },
  ];

  choirs.forEach((choir, idx) => {
    const choirNum = idx + 1;

    // Big Bead: Salutation
    steps.push({
      id: `choir-${choirNum}-salutation`,
      title: `${choir.name}`,
      subtitle: `Virtue: ${choir.virtue}`,
      prayerId: choir.prayerId,
      beadText: `Big Bead: ${choir.name}`,
      choirIndex: choirNum,
      choirName: choir.name,
      intention: choir.virtue,
      section: 'choirs'
    });

    // Big Bead: 1 Our Father
    steps.push({
      id: `choir-${choirNum}-our-father`,
      title: `${choir.name} - Our Father`,
      subtitle: `1 Our Father for ${choir.name}`,
      prayerId: 'our-father',
      beadText: `Big Bead: 1 Our Father`,
      choirIndex: choirNum,
      choirName: choir.name,
      section: 'choirs'
    });

    // 3 Small Beads: 3 Hail Marys
    for (let h = 1; h <= 3; h++) {
      steps.push({
        id: `choir-${choirNum}-hail-mary-${h}`,
        title: `${choir.name} - Hail Mary ${h}/3`,
        subtitle: `Small Bead ${h} of 3`,
        prayerId: 'hail-mary',
        beadText: `Small Bead ${h}/3: Hail Mary`,
        choirIndex: choirNum,
        choirName: choir.name,
        section: 'choirs'
      });
    }
  });

  // The 4 Pendant Beads (Pendant Tail)
  const pendantOurFathers = [
    { title: 'In Honor of St. Michael', honor: 'St. Michael the Archangel' },
    { title: 'In Honor of St. Gabriel', honor: 'St. Gabriel the Archangel' },
    { title: 'In Honor of St. Raphael', honor: 'St. Raphael the Archangel' },
    { title: 'In Honor of Our Guardian Angel', honor: 'Our Guardian Angel' }
  ];

  pendantOurFathers.forEach((item, index) => {
    steps.push({
      id: `pendant-our-father-${index + 1}`,
      title: `Our Father ${index + 1}/4 (${item.honor})`,
      subtitle: `Pendant Bead ${index + 1} of 4: ${item.title}`,
      prayerId: 'our-father',
      beadText: `Pendant Bead ${index + 1}: ${item.honor}`,
      section: 'pendant'
    });
  });

  // Concluding Prayers & Sign of the Cross
  steps.push({
    id: 'st-michael-concluding-invocation',
    title: 'O Glorious Prince St. Michael',
    subtitle: 'Concluding Invocations & Versicle',
    prayerId: 'st-michael-concluding-prayer',
    section: 'concluding'
  });

  steps.push({
    id: 'st-michael-main-prayer',
    title: 'Prayer to St. Michael the Archangel (Leo XIII)',
    subtitle: 'Prince of the Heavenly Host',
    prayerId: 'st-michael-closing-prayer',
    section: 'concluding'
  });

  steps.push({
    id: 'concluding-sign-of-the-cross',
    title: 'Sign of the Cross (The End)',
    subtitle: 'In the name of the Father, and of the Son, and of the Holy Spirit. Amen.',
    prayerId: 'sign-of-the-cross',
    section: 'concluding'
  });

  return steps;
}
