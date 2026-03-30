export type InterlinearPair = { la: string; en: string };

/**
 * Global dictionary of Latin words/short phrases and their English translations.
 * Key: Latin (case-sensitive, includes punctuation if needed to match source text)
 * Value: English translation
 */
const DICTIONARY: Record<string, string> = {
    // Shared / Particles
    'Amen': 'Amen',
    'et': 'and',
    'In': 'In',
    'in': 'in',
    'Et': 'And',
    'O': 'O',
    'o': 'o',
    'qui': 'who',
    'est': 'is',
    'sed': 'but',
    'pro': 'for',
    'ad': 'to',
    'Ad': 'To',
    'ab': 'from',
    'ex': 'from',
    'de': 'by',
    'tu': 'thou',
    'tuum': 'Thy',
    'tua': 'Thy',
    'tui': 'thy',
    'tuae': 'Your',
    'Te': 'You',
    'te': 'thee',
    'nos': 'us',
    'nobis': 'us',
    'nostra': 'our',
    'nostrae': 'our',
    'nostrum': 'our',
    'noster': 'our',
    'eius': 'His',
    'quae': 'who',
    'ut': 'as',
    'sicut': 'as',

    // Sign of the Cross
    'nomine': 'the name',
    'Patris': 'of the Father',
    'Spiritus': 'of the Spirit',
    'Sancti': 'Holy',

    // Divine Mercy Opening prayers
    'Expirasti': 'You expired',
    'Iesu': 'Jesus',
    'fons': 'the source',
    'vitae': 'of life',
    'animabus': 'souls',
    'emanavit': 'gushed forth',
    'oceanus': 'the ocean',
    'misericordiae': 'of mercy',
    'toto': 'the whole',
    'mundo': 'world',
    'apertus': 'opened up',
    'Fons': 'Fount',
    'Vitae': 'of Life',
    'inscrutabilis': 'unfathomable',
    'Divina': 'Divine',
    'Misericordia': 'Mercy',
    'totum': 'the whole',
    'mundum': 'world',
    'circumda': 'envelop',
    'teipsum': 'empty Yourself',
    'super': 'upon',
    'effunde': 'out',
    'Sanguis': 'Blood',
    'Aqua': 'Water',
    'corde': 'the Heart',
    'confido': 'I trust!',
    'emanasti': 'gushed forth',

    // Our Father
    'Pater': 'Father',
    'es': 'art',
    'caelis': 'heaven',
    'sanctificetur': 'hallowed be',
    'nomen': 'name',
    'Adveniat': 'come',
    'regnum': 'kingdom',
    'Fiat': 'be done',
    'voluntas': 'will',
    'caelo': 'heaven',
    'terra': 'earth',
    'Panem': 'bread',
    'quotidianum': 'daily',
    'da': 'give',
    'hodie': 'this day',
    'dimitte': 'forgive',
    'debita': 'trespasses',
    'dimittimus': 'forgive',
    'debitoribus': 'those who trespass',
    'ne': 'not',
    'inducas': 'lead',
    'tentationem': 'temptation',
    'libera': 'deliver',
    'malo': 'evil',
    'omnipotentem': 'almighty',
    'omnipotentis': 'almighty',

    // Hail Mary
    'Ave': 'Hail',
    'Maria': 'Mary',
    'gratia': 'of grace',
    'plena': 'full',
    'Dominus': 'the Lord',
    'tecum': 'is with thee',
    'Benedicta': 'Blessed',
    'mulieribus': 'among women',
    'benedictus': 'blessed',
    'fructus': 'is the fruit',
    'ventris': 'of the womb',
    'Iesus': 'Jesus',
    'Sancta': 'Holy',
    'Mater': 'Mother',
    'Dei': 'of God',
    'ora': 'pray',
    'peccatoribus': 'sinners',
    'nunc': 'now',
    'hora': 'the hour',
    'mortis': 'of death',

    // Glory Be
    'Gloria': 'Glory',
    'Patri': 'to the Father',
    'Filio': 'to the Son',
    'Spiritui': 'to the Spirit',
    'Sancto': 'Holy',
    'erat': 'it was',
    'principio': 'the beginning',
    'semper': 'ever',
    'saecula': 'without',
    'saeculorum': 'end',

    // Apostles' Creed
    'Credo': 'I believe',
    'Deum': 'God',
    'Patrem': 'the Father',
    'Creatorem': 'Creator',
    'caeli': 'of heaven',
    'terrae': 'earth',
    'Iesum': 'Jesus',
    'Christum': 'Christ',
    'Filium': 'Son',
    'unicum': 'only',
    'Dominum': 'Lord',
    'conceptus': 'was conceived',
    'Spiritu': 'the Spirit',
    'natus': 'born',
    'Virgine': 'the Virgin',
    'passus': 'suffered',
    'sub': 'under',
    'Pontio': 'Pontius',
    'Pilato': 'Pilate',
    'crucifixus': 'was crucified',
    'mortuus': 'died',
    'sepultus': 'was buried',
    'descendit': 'He descended',
    'inferos': 'hell',
    'tertia': 'on the third',
    'die': 'day',
    'resurrexit': 'He rose again',
    'a': 'from',
    'mortuis': 'the dead',
    'ascendit': 'He ascended',
    'sedet': 'is seated',
    'dexteram': 'the right hand',
    'inde': 'from there',
    'venturus': 'He will come',
    'iudicare': 'to judge',
    'vivos': 'the living',
    'mortuos': 'the dead',
    'Spiritum': 'the Spirit',
    'Sanctum': 'Holy',
    'sanctam': 'the holy',
    'Ecclesiam': 'Church',
    'catholicam': 'catholic',
    'sanctorum': 'of saints',
    'communionem': 'the communion',
    'remissionem': 'the forgiveness',
    'peccatorum': 'of sins',
    'carnis': 'of the body',
    'resurrectionem': 'the resurrection',
    'vitam': 'life',
    'aeternam': 'everlasting',

    // Eternal Father
    'Aeterne': 'Eternal',
    'offero': 'I offer',
    'Tibi': 'to You',
    'Corpus': 'the Body',
    'Sanguinem': 'Blood',
    'Animam': 'Soul',
    'Divinitatem': 'Divinity',
    'delectissimi': 'of Your dearly beloved',
    'Filii': 'Son',
    'Tui': 'Your',
    'Domini': 'Lord',
    'Christi': 'Christ',
    'propitiatione': 'atonement',
    'peccatis': 'sins',
    'nostris': 'our',
    'totius': 'of the whole',
    'mundi': 'world',

    // Sorrowful Passion
    'Propter': 'For the sake of',
    'sorrowful': 'sorrowful',
    'Passionem': 'Passion',
    'miserere': 'have mercy',

    // Holy God
    'Sanctus': 'Holy',
    'Deus': 'God',
    'Fortis': 'Mighty One',
    'Immortalis': 'Immortal One',

    // Fatima Prayer
    'mi': 'my',
    'igne': 'the fire',
    'inferni': 'of hell',
    'perduc': 'lead',
    'caelum': 'heaven',
    'omnes': 'all',
    'animas': 'souls',
    'praesertim': 'especially',
    'eas': 'those',
    'maxime': 'most',
    'indigent': 'have need',

    // Hail Holy Queen
    'Salve': 'Hail',
    'Regina': 'Queen',
    'mater': 'mother',
    'vita': 'our life',
    'dulcedo': 'our sweetness',
    'spes': 'hope',
    'salve': 'hail',
    'clamamus': 'do we cry',
    'exsules': 'poor banished',
    'filii': 'children',
    'Hevae': 'of Eve',
    'suspiramus': 'do we send up our sighs',
    'gementes': 'mourning',
    'flentes': 'weeping',
    'hac': 'this',
    'lacrimarum': 'of tears',
    'valle': 'valley',

    // Special cases
    'gratia plena': 'full of grace'
};


/**
 * Generates an array of interlinear pairs by looking up words/phrases in the dictionary.
 * @param latin The Latin text to translate.
 * @returns An array of Latin-English pairs.
 */
export const getInterlinearPairs = (latin?: string): InterlinearPair[] | null => {
    if (!latin) return null;

    // Normalizing whitespace and splitting
    const words = latin.trim().split(/\s+/);

    // Check if we have at least one valid translation to avoid returning an empty/useless list
    let hasTranslations = false;

    const result = words.map(la => {
        // Strip punctuation (.,;:!?"') for lookup
        const cleanLa = la.replace(/[.,;:!?"']/g, '');

        // Try exact match first
        let en = DICTIONARY[cleanLa] || '';

        // If not found, try lowercase
        if (!en && cleanLa.length > 1) {
            en = DICTIONARY[cleanLa.toLowerCase()] || '';
        }

        if (en) hasTranslations = true;
        return { la, en };
    });

    return hasTranslations ? result : null;
};
