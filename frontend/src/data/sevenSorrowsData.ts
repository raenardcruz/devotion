export interface SorrowVerse {
  text: string;
  ref: string;
}

export interface SevenSorrowMystery {
  index: number;
  ordinal: string;
  latinOrdinal: string;
  title: string;
  shortTitle: string;
  latinTitle: string;
  description: string;
  meditation: string;
  virtues: string[];
  verse: SorrowVerse;
  beadVerses: SorrowVerse[];
  image: string;
  artist: string;
}

export type SevenSorrowsStepType =
  | 'opening'
  | 'intentions'
  | 'intro'
  | 'sorrow-header'
  | 'sorrow-our-father'
  | 'sorrow-hail-mary'
  | 'sorrow-invocation'
  | 'tears-hail-mary'
  | 'closing'
  | 'completed';

export interface SevenSorrowsStep {
  id: string;
  type: SevenSorrowsStepType;
  section: 'intro' | 'sorrows' | 'tears' | 'concluding';
  title: string;
  latinTitle?: string;
  subtitle?: string;
  prayerId: string;
  sorrowIndex?: number;
  sorrowTitle?: string;
  latinSorrowTitle?: string;
  beadNumber?: number; // 0 for Our Father, 1..7 for Hail Marys in a sorrow, or 1..3 for Tears
  beadText?: string;
  verse?: SorrowVerse;
  mysteryVerse?: SorrowVerse;
  description?: string;
  virtues?: string[];
  customText?: string;
  customLatin?: string;
}

export const SEVEN_SORROWS_PROMISES = [
  {
    title: 'Peace in Families',
    desc: 'I will grant peace to their families and enlighten them about the divine mysteries.'
  },
  {
    title: 'Consolation in Pains',
    desc: 'I will console them in their pains and accompany them in their work.'
  },
  {
    title: 'Defense in Spiritual Battles',
    desc: 'I will defend them in their spiritual battles with the infernal enemy and protect them at every instant of their lives.'
  },
  {
    title: 'Assistance at the Hour of Death',
    desc: 'I will visibly help them at the moment of their death—they will see the face of their Mother.'
  }
];

export const SEVEN_SORROWS_MYSTERIES: SevenSorrowMystery[] = [
  {
    index: 1,
    ordinal: 'First Sorrow',
    latinOrdinal: 'Dolor Primus',
    title: 'The Prophecy of Simeon',
    shortTitle: '1. Simeon',
    latinTitle: 'Prophetia Simeonis',
    description:
      'When Mary and Joseph present the Infant Jesus in the Temple, the righteous Simeon foretells that Jesus will be a sign of contradiction and that a sword of sorrow will pierce Mary’s own soul.',
    meditation:
      'Gaze upon Our Lady in the Temple as she hears Simeon’s prophecy. Though her heart is pierced with the foreknowledge of Calvary, she surrenders completely to the Father’s salvific plan.',
    virtues: ['Submission to God’s Will', 'Trust in Divine Providence'],
    image: '/images/seven-sorrows/sorrow_1_simeon.jpg',
    artist: 'Rembrandt van Rijn — Simeon in the Temple (c. 1669)',
    verse: {
      text: 'And Simeon blessed them and said to Mary his mother, "Behold, this child is set for the fall and rising of many in Israel, and for a sign that is spoken against (and a sword will pierce through your own soul also), that thoughts out of many hearts may be revealed."',
      ref: 'Luke 2:34-35'
    },
    beadVerses: [
      { text: 'When the time came for their purification according to the law of Moses, they brought him up to Jerusalem to present him to the Lord.', ref: 'Luke 2:22' },
      { text: 'Now there was a man in Jerusalem, whose name was Simeon, and this man was righteous and devout, looking for the consolation of Israel.', ref: 'Luke 2:25' },
      { text: 'Inspired by the Spirit he came into the temple; and when the parents brought in the child Jesus, he took him up in his arms and blessed God.', ref: 'Luke 2:27-28' },
      { text: '"Lord, now lettest thou thy servant depart in peace, according to thy word; for mine eyes have seen thy salvation."', ref: 'Luke 2:29-30' },
      { text: 'And his father and his mother marveled at what was said about him.', ref: 'Luke 2:33' },
      { text: 'And Simeon blessed them and said to Mary his mother, "Behold, this child is set for the fall and rising of many in Israel, and for a sign that is spoken against."', ref: 'Luke 2:34' },
      { text: '"And a sword will pierce through your own soul also, that thoughts out of many hearts may be revealed."', ref: 'Luke 2:35' }
    ]
  },
  {
    index: 2,
    ordinal: 'Second Sorrow',
    latinOrdinal: 'Dolor Secundus',
    title: 'The Flight into Egypt',
    shortTitle: '2. Egypt',
    latinTitle: 'Fuga in Aegyptum',
    description:
      'Warned by an angel in the night that King Herod seeks to destroy the Child Jesus, Mary and Joseph flee with the newborn Savior into exile in a foreign land.',
    meditation:
      'Walk beside Mary and Joseph through the dark desert toward Egypt. Pray for all who endure exile, persecution, and uncertainty, and ask for the grace to protect the presence of Christ in your soul.',
    virtues: ['Courage in Exile', 'Prompt Obedience'],
    image: '/images/seven-sorrows/sorrow_2_egypt.jpg',
    artist: 'Bartolomé Esteban Murillo — The Flight into Egypt (c. 1647)',
    verse: {
      text: 'Now when they had departed, behold, an angel of the Lord appeared to Joseph in a dream and said, "Rise, take the child and his mother, and flee to Egypt, and remain there till I tell you; for Herod is about to search for the child, to destroy him." And he rose and took the child and his mother by night, and departed to Egypt.',
      ref: 'Matthew 2:13-14'
    },
    beadVerses: [
      { text: 'Behold, an angel of the Lord appeared to Joseph in a dream and said, "Rise, take the child and his mother, and flee to Egypt."', ref: 'Matthew 2:13a' },
      { text: '"And remain there till I tell you; for Herod is about to search for the child, to destroy him."', ref: 'Matthew 2:13b' },
      { text: 'And he rose and took the child and his mother by night, and departed to Egypt.', ref: 'Matthew 2:14' },
      { text: 'And remained there until the death of Herod. This was to fulfil what the Lord had spoken by the prophet, "Out of Egypt have I called my son."', ref: 'Matthew 2:15' },
      { text: 'Then Herod, when he saw that he had been tricked by the wise men, was in a furious rage.', ref: 'Matthew 2:16' },
      { text: '"A voice was heard in Ramah, wailing and loud lamentation, Rachel weeping for her children."', ref: 'Matthew 2:18' },
      { text: 'The Lord is my light and my salvation; whom shall I fear? The Lord is the stronghold of my life.', ref: 'Psalm 27:1' }
    ]
  },
  {
    index: 3,
    ordinal: 'Third Sorrow',
    latinOrdinal: 'Dolor Tertius',
    title: 'The Loss of the Child Jesus in the Temple',
    shortTitle: '3. Temple',
    latinTitle: 'Amissio Pueri Iesu in Templo',
    description:
      'Returning from the Passover feast in Jerusalem, Mary and Joseph discover that twelve-year-old Jesus is missing and search for Him in deep anguish for three days.',
    meditation:
      'Share in the three days of sorrow Mary endured while separated from Jesus. Whenever sin or spiritual dryness obscures His face, ask Our Lady for the perseverance to seek Him above all things.',
    virtues: ['Zeal for Jesus', 'Repentance for Losing Grace'],
    image: '/images/seven-sorrows/sorrow_3_temple.jpg',
    artist: 'Abraham Hondius — Christ among the Doctors (1668)',
    verse: {
      text: 'And when the feast was ended, as they were returning, the boy Jesus stayed behind in Jerusalem. His parents did not know it... and when they did not find him, they returned to Jerusalem, seeking him. And his mother said to him, "Son, why have you treated us so? Behold, your father and I have been looking for you anxiously."',
      ref: 'Luke 2:43-48'
    },
    beadVerses: [
      { text: 'And when he was twelve years old, they went up to Jerusalem according to the custom of the feast.', ref: 'Luke 2:42' },
      { text: 'When the feast was ended, as they were returning, the boy Jesus stayed behind in Jerusalem. His parents did not know it.', ref: 'Luke 2:43' },
      { text: 'Supposing him to be in the company they went a day’s journey, and they sought him among their kinsfolk and acquaintances.', ref: 'Luke 2:44' },
      { text: 'And when they did not find him, they returned to Jerusalem, seeking him.', ref: 'Luke 2:45' },
      { text: 'After three days they found him in the temple, sitting among the teachers, listening to them and asking them questions.', ref: 'Luke 2:46' },
      { text: 'His mother said to him, "Son, why have you treated us so? Behold, your father and I have been looking for you anxiously."', ref: 'Luke 2:48' },
      { text: 'And he went down with them and came to Nazareth, and was obedient to them; and his mother kept all these things in her heart.', ref: 'Luke 2:51' }
    ]
  },
  {
    index: 4,
    ordinal: 'Fourth Sorrow',
    latinOrdinal: 'Dolor Quartus',
    title: 'Mary Meets Jesus on the Way to Calvary',
    shortTitle: '4. Calvary',
    latinTitle: 'Occursus Mariae et Iesu in Via Crucis',
    description:
      'On the sorrowful way to Golgotha, Mary meets her beloved Son bruised, crowned with thorns, and crushed beneath the weight of the Cross.',
    meditation:
      'Behold the silent, piercing gaze exchanged between Mother and Son on the Via Dolorosa. Unite your daily crosses with theirs, never turning away from those who suffer.',
    virtues: ['Patience in Suffering', 'Compassionate Love'],
    image: '/images/seven-sorrows/sorrow_4_calvary.jpg',
    artist: 'El Greco — Christ Carrying the Cross (c. 1580)',
    verse: {
      text: 'And there followed him a great multitude of the people, and of women who bewailed and lamented him. But Jesus turning to them said, "Daughters of Jerusalem, do not weep for me, but weep for yourselves and for your children."',
      ref: 'Luke 23:27-28'
    },
    beadVerses: [
      { text: 'So they took Jesus, and he went out, bearing his own cross, to the place called the place of a skull, which is called in Hebrew Golgotha.', ref: 'John 19:17' },
      { text: 'And as they led him away, they seized one Simon of Cyrene, who was coming in from the country, and laid on him the cross, to carry it behind Jesus.', ref: 'Luke 23:26' },
      { text: 'And there followed him a great multitude of the people, and of women who bewailed and lamented him.', ref: 'Luke 23:27' },
      { text: 'But Jesus turning to them said, "Daughters of Jerusalem, do not weep for me, but weep for yourselves and for your children."', ref: 'Luke 23:28' },
      { text: 'He was despised and rejected by men; a man of sorrows, and acquainted with grief.', ref: 'Isaiah 53:3' },
      { text: 'Surely he has borne our griefs and carried our sorrows.', ref: 'Isaiah 53:4' },
      { text: '"Is it nothing to you, all you who pass by? Look and see if there is any sorrow like my sorrow."', ref: 'Lamentations 1:12' }
    ]
  },
  {
    index: 5,
    ordinal: 'Fifth Sorrow',
    latinOrdinal: 'Dolor Quintus',
    title: 'The Crucifixion and Death of Jesus',
    shortTitle: '5. Crucifixion',
    latinTitle: 'Crucifixio et Mors Iesu',
    description:
      'Mary stands courageously beneath the Cross as Jesus is nailed to the wood, endures three hours of agony, entrusts us to her maternal care, and breathes His last.',
    meditation:
      'Stand with Mary at the foot of the Cross—Stabat Mater Dolorosa. Hear Jesus say, "Behold, your mother," and welcome her into the innermost sanctuary of your life.',
    virtues: ['Self-Offering with Christ', 'Fidelity in Trial'],
    image: '/images/seven-sorrows/sorrow_5_crucifixion.jpg',
    artist: 'Hendrick ter Brugghen — The Crucifixion with the Virgin and St. John (c. 1625)',
    verse: {
      text: 'Standing by the cross of Jesus were his mother, and his mother’s sister, Mary the wife of Clopas, and Mary Magdalene. When Jesus saw his mother, and the disciple whom he loved standing near, he said to his mother, "Woman, behold, your son!" Then he said to the disciple, "Behold, your mother!" And from that hour the disciple took her to his own home.',
      ref: 'John 19:25-27'
    },
    beadVerses: [
      { text: 'And when they came to the place which is called The Skull, there they crucified him.', ref: 'Luke 23:33' },
      { text: 'And Jesus said, "Father, forgive them; for they know not what they do."', ref: 'Luke 23:34' },
      { text: 'Standing by the cross of Jesus were his mother, and his mother’s sister, Mary the wife of Clopas, and Mary Magdalene.', ref: 'John 19:25' },
      { text: 'When Jesus saw his mother, and the disciple whom he loved standing near, he said to his mother, "Woman, behold, your son!"', ref: 'John 19:26' },
      { text: 'Then he said to the disciple, "Behold, your mother!" And from that hour the disciple took her to his own home.', ref: 'John 19:27' },
      { text: 'When Jesus had received the vinegar, he said, "It is finished"; and he bowed his head and gave up his spirit.', ref: 'John 19:30' },
      { text: 'Then Jesus, crying with a loud voice, said, "Father, into thy hands I commit my spirit!" And having said this he breathed his last.', ref: 'Luke 23:46' }
    ]
  },
  {
    index: 6,
    ordinal: 'Sixth Sorrow',
    latinOrdinal: 'Dolor Sextus',
    title: 'Jesus is Pierced and Taken Down from the Cross',
    shortTitle: '6. Pietà',
    latinTitle: 'Lanceatio et Depositio Iesu de Cruce',
    description:
      'A soldier pierces the Sacred Heart of Jesus with a lance, and His lifeless Body is taken down from the Cross and laid in the arms of His Sorrowful Mother.',
    meditation:
      'Contemplate the Pietà—Mary cradling the lifeless Body of the Son she once held in Bethlehem. May every wound of Christ be engraved upon our hearts with gratitude and contrition.',
    virtues: ['Deep Contrition for Sin', 'Tender Devotion to the Passion'],
    image: '/images/seven-sorrows/sorrow_6_pieta.jpg',
    artist: 'William-Adolphe Bouguereau — Pietà (1876)',
    verse: {
      text: 'But when they came to Jesus and saw that he was already dead, they did not break his legs. But one of the soldiers pierced his side with a spear, and at once there came out blood and water... And Joseph of Arimathea took the body of Jesus down from the cross.',
      ref: 'John 19:33-34, 38'
    },
    beadVerses: [
      { text: 'When they came to Jesus and saw that he was already dead, they did not break his legs.', ref: 'John 19:33' },
      { text: 'But one of the soldiers pierced his side with a spear, and at once there came out blood and water.', ref: 'John 19:34' },
      { text: 'For these things took place that the scripture might be fulfilled, "Not a bone of him shall be broken."', ref: 'John 19:36' },
      { text: 'And again another scripture says, "They shall look on him whom they have pierced."', ref: 'John 19:37' },
      { text: 'After this Joseph of Arimathea, who was a disciple of Jesus, asked Pilate that he might take away the body of Jesus.', ref: 'John 19:38' },
      { text: 'Then he took it down from the cross and placed Him in the arms of His afflicted Mother.', ref: 'Luke 23:53a' },
      { text: '"O all you who pass by the way, attend and see if there be any sorrow like to my sorrow."', ref: 'Lamentations 1:12' }
    ]
  },
  {
    index: 7,
    ordinal: 'Seventh Sorrow',
    latinOrdinal: 'Dolor Septimus',
    title: 'The Burial of Jesus in the Tomb',
    shortTitle: '7. Burial',
    latinTitle: 'Sepultura Iesu',
    description:
      'The sacred Body of Jesus is wrapped in linen with spices and laid in the garden tomb, and the heavy stone is rolled across the entrance while Mary awaits the dawn of the Resurrection.',
    meditation:
      'Stand with Our Lady beside the sealed tomb in the quiet of Holy Saturday. Even in the deepest silence and grief, her unshakeable faith holds fast to the promise of the Resurrection.',
    virtues: ['Unwavering Hope', 'Perseverance in Faith'],
    image: '/images/seven-sorrows/sorrow_7_burial.jpg',
    artist: 'Caravaggio — The Entombment of Christ (c. 1603)',
    verse: {
      text: 'They took the body of Jesus, and bound it in linen cloths with the spices, as is the burial custom of the Jews. Now in the place where he was crucified there was a garden, and in the garden a new tomb where no one had ever been laid. So because of the Jewish day of Preparation, as the tomb was close at hand, they laid Jesus there.',
      ref: 'John 19:40-42'
    },
    beadVerses: [
      { text: 'Nicodemus also, who had at first come to him by night, came bringing a mixture of myrrh and aloes, about a hundred pounds’ weight.', ref: 'John 19:39' },
      { text: 'They took the body of Jesus, and bound it in linen cloths with the spices, as is the burial custom of the Jews.', ref: 'John 19:40' },
      { text: 'Now in the place where he was crucified there was a garden, and in the garden a new tomb where no one had ever been laid.', ref: 'John 19:41' },
      { text: 'So because of the Jewish day of Preparation, as the tomb was close at hand, they laid Jesus there.', ref: 'John 19:42' },
      { text: 'And Joseph rolled a great stone to the door of the tomb, and departed.', ref: 'Matthew 27:60' },
      { text: 'The women who had come with him from Galilee followed, and saw the tomb, and how his body was laid.', ref: 'Luke 23:55' },
      { text: 'Unless a grain of wheat falls into the earth and dies, it remains alone; but if it dies, it bears much fruit.', ref: 'John 12:24' }
    ]
  }
];

export function generateSevenSorrowsSteps(hasIntentions: boolean = false): SevenSorrowsStep[] {
  const steps: SevenSorrowsStep[] = [];

  // 1. Opening Sign of the Cross
  steps.push({
    id: 'opening-sign-of-the-cross',
    type: 'opening',
    section: 'intro',
    title: 'The Sign of the Cross',
    latinTitle: 'Signum Crucis',
    subtitle: 'In the name of the Father, and of the Son, and of the Holy Spirit',
    prayerId: 'sign-of-the-cross'
  });

  // 2. Prayer Intentions (if user has active intentions)
  if (hasIntentions) {
    steps.push({
      id: 'prayer-intentions',
      type: 'intentions',
      section: 'intro',
      title: 'Prayer Intentions',
      latinTitle: 'Intentiones Orationis',
      subtitle: 'Petitions entrusted to Our Lady of Sorrows',
      prayerId: 'prayer-intentions'
    });
  }

  // 3. Introductory Offering Prayer
  steps.push({
    id: 'intro-offering',
    type: 'intro',
    section: 'intro',
    title: 'Introductory Offering',
    latinTitle: 'Oblatio Initialis Septem Dolorum',
    subtitle: 'Offering the Rosary of the Seven Sorrows',
    prayerId: 'seven-sorrows-opening'
  });

  // 4. Act of Contrition
  steps.push({
    id: 'intro-act-of-contrition',
    type: 'intro',
    section: 'intro',
    title: 'Act of Contrition',
    latinTitle: 'Actus Contritionis',
    subtitle: 'True Sorrow for Our Sins',
    prayerId: 'act-of-contrition'
  });

  // 5. The 7 Sorrows (Mysteries)
  SEVEN_SORROWS_MYSTERIES.forEach((mystery) => {
    const sIdx = mystery.index;

    // Step A: Mystery Announcement + Bible Verse (Displayed BEFORE starting each mystery)
    steps.push({
      id: `sorrow-${sIdx}-header`,
      type: 'sorrow-header',
      section: 'sorrows',
      title: `${mystery.ordinal}: ${mystery.title}`,
      latinTitle: `${mystery.latinOrdinal}: ${mystery.latinTitle}`,
      subtitle: 'Scripture Reading & Meditation Before Starting Mystery',
      prayerId: '',
      sorrowIndex: sIdx,
      sorrowTitle: mystery.title,
      latinSorrowTitle: mystery.latinTitle,
      verse: mystery.verse,
      mysteryVerse: mystery.verse,
      description: mystery.description,
      virtues: mystery.virtues,
      beadText: `${mystery.ordinal} — Scripture`
    });

    // Step B: 1 Our Father
    steps.push({
      id: `sorrow-${sIdx}-our-father`,
      type: 'sorrow-our-father',
      section: 'sorrows',
      title: 'Our Father',
      latinTitle: 'Pater Noster',
      subtitle: `${mystery.ordinal}: ${mystery.title}`,
      prayerId: 'our-father',
      sorrowIndex: sIdx,
      sorrowTitle: mystery.title,
      latinSorrowTitle: mystery.latinTitle,
      beadNumber: 0,
      beadText: 'Large Medal Bead: 1 Our Father',
      mysteryVerse: mystery.verse
    });

    // Step C: 7 Hail Marys
    for (let h = 1; h <= 7; h++) {
      steps.push({
        id: `sorrow-${sIdx}-hail-mary-${h}`,
        type: 'sorrow-hail-mary',
        section: 'sorrows',
        title: `Hail Mary (${h} of 7)`,
        latinTitle: `Ave Maria (${h} / 7)`,
        subtitle: `${mystery.ordinal}: ${mystery.title}`,
        prayerId: 'hail-mary',
        sorrowIndex: sIdx,
        sorrowTitle: mystery.title,
        latinSorrowTitle: mystery.latinTitle,
        beadNumber: h,
        beadText: `Sorrow Bead ${h} of 7`,
        verse: mystery.beadVerses[h - 1] || mystery.verse,
        mysteryVerse: mystery.verse
      });
    }

    // Step D: Marian Sorrow Invocation
    steps.push({
      id: `sorrow-${sIdx}-invocation`,
      type: 'sorrow-invocation',
      section: 'sorrows',
      title: 'Invocation to Our Sorrowful Mother',
      latinTitle: 'Invocatio Matris Dolorosae',
      subtitle: `Concluding the ${mystery.ordinal}`,
      prayerId: 'seven-sorrows-invocation',
      sorrowIndex: sIdx,
      sorrowTitle: mystery.title,
      latinSorrowTitle: mystery.latinTitle,
      beadText: 'Closing Invocation',
      mysteryVerse: mystery.verse
    });
  });

  // 6. Three Concluding Hail Marys in Honor of the Tears of Our Lady
  for (let t = 1; t <= 3; t++) {
    steps.push({
      id: `tears-hail-mary-${t}`,
      type: 'tears-hail-mary',
      section: 'tears',
      title: `Hail Mary (${t} of 3) — Tears of Mary`,
      latinTitle: `Ave Maria (${t} / 3) — In Honorem Lacrimarum Mariae`,
      subtitle: 'In honor of the tears shed by Our Sorrowful Mother',
      prayerId: 'hail-mary',
      beadNumber: t,
      beadText: `Tears Bead ${t} of 3`
    });
  }

  // 7. Concluding Prayer of the Seven Sorrows
  steps.push({
    id: 'closing-seven-sorrows-prayer',
    type: 'closing',
    section: 'concluding',
    title: 'Concluding Prayer of the Seven Sorrows',
    latinTitle: 'Oratio Conclusiva Septem Dolorum',
    subtitle: 'Liturgical Collect of Our Lady of Sorrows',
    prayerId: 'seven-sorrows-closing'
  });

  // 8. Final Sign of the Cross
  steps.push({
    id: 'closing-sign-of-the-cross',
    type: 'closing',
    section: 'concluding',
    title: 'The Sign of the Cross',
    latinTitle: 'Signum Crucis',
    subtitle: 'In the name of the Father, and of the Son, and of the Holy Spirit. Amen.',
    prayerId: 'sign-of-the-cross'
  });

  // 9. Completion Step
  steps.push({
    id: 'seven-sorrows-completed',
    type: 'completed',
    section: 'concluding',
    title: 'Seven Sorrows Chaplet Completed',
    latinTitle: 'Corolla Septem Dolorum Finita Est',
    prayerId: 'completion-marker',
    customText: 'You have completed the Chaplet of the Seven Sorrows. May Our Lady of Sorrows console your heart and lead you ever closer to her Son, Jesus.',
    customLatin: 'Corollam Septem Dolorum complevisti. Beata Virgo Maria Dolorosa cor tuum consoletur et ad Filium suum Iesum te perducat.'
  });

  return steps;
}
