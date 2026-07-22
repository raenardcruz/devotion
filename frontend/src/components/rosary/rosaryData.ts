export type Verse = {
    text: string;
    ref: string;
};

export type Mystery = {
    title: string;
    description: string;
    virtues: string[];
    verses: Verse[];
    image: string;
};

export type RosaryData = {
    [key: string]: Mystery[];
};

export type RosaryStepType = 'intro' | 'mystery-header' | 'decade-bead' | 'decade-start' | 'decade-end' | 'closing' | 'opening';

export type RosaryStep = {
    type: RosaryStepType;
    prayerId: string;
    title?: string;
    content?: string;
    latin?: string;
    verse?: Verse;
    mysteryTitle?: string;
    beadNumber?: number;
    decadeNumber?: number;
    youtube?: string;
    description?: string;
    virtues?: string[];
};

export const ROSARY_DATA: RosaryData = {
    "Joyful": [
        {
            "title": "The Annunciation",
            "description": "The Angel Gabriel announces to Mary that she will conceive the Son of God.",
            "virtues": ["Humility"],
            "image": "/images/rosary/annunciation.png",
            "verses": [
                {
                    "text": "In the sixth month the angel Gabriel was sent from God to a city of Galilee named Nazareth,",
                    "ref": "Luke 1:26"
                },
                {
                    "text": "to a virgin betrothed to a man whose name was Joseph, of the house of David; and the virgin's name was Mary.",
                    "ref": "Luke 1:27"
                },
                {
                    "text": "And he came to her and said, \"Hail, O favored one, the Lord is with you!\"",
                    "ref": "Luke 1:28"
                },
                {
                    "text": "But she was greatly troubled at the saying, and considered in her mind what sort of greeting this might be.",
                    "ref": "Luke 1:29"
                },
                {
                    "text": "And the angel said to her, \"Do not be afraid, Mary, for you have found favor with God.",
                    "ref": "Luke 1:30"
                },
                {
                    "text": "And behold, you will conceive in your womb and bear a son, and you shall call his name Jesus.",
                    "ref": "Luke 1:31"
                },
                {
                    "text": "He will be great, and will be called the Son of the Most High; and the Lord God will give to him the throne of his father David,",
                    "ref": "Luke 1:32"
                },
                {
                    "text": "And Mary said to the angel, \"How shall this be, since I have no husband?\"",
                    "ref": "Luke 1:34"
                },
                {
                    "text": "And the angel said to her, \"The Holy Spirit will come upon you, and the power of the Most High will overshadow you; therefore the child to be born will be called holy, the Son of God.",
                    "ref": "Luke 1:35"
                },
                {
                    "text": "And Mary said, \"Behold, I am the handmaid of the Lord; let it be to me according to your word.\" And the angel departed from her.",
                    "ref": "Luke 1:38"
                }
            ]
        },
        {
            "title": "The Visitation",
            "description": "Mary visits her cousin Elizabeth, who is pregnant with John the Baptist.",
            "virtues": ["Charity", "Love of Neighbor"],
            "image": "/images/rosary/visitation.png",
            "verses": [
                {
                    "text": "In those days Mary arose and went with haste into the hill country, to a city of Judah,",
                    "ref": "Luke 1:39"
                },
                {
                    "text": "and she entered the house of Zechariah and greeted Elizabeth.",
                    "ref": "Luke 1:40"
                },
                {
                    "text": "And when Elizabeth heard the greeting of Mary, the babe leaped in her womb; and Elizabeth was filled with the Holy Spirit",
                    "ref": "Luke 1:41"
                },
                {
                    "text": "and she exclaimed with a loud cry, \"Blessed are you among women, and blessed is the fruit of your womb!",
                    "ref": "Luke 1:42"
                },
                {
                    "text": "And why is this granted me, that the mother of my Lord should come to me?",
                    "ref": "Luke 1:43"
                },
                {
                    "text": "For behold, when the voice of your greeting came to my ears, the babe in my womb leaped for joy.",
                    "ref": "Luke 1:44"
                },
                {
                    "text": "And blessed is she who believed that there would be a fulfilment of what was spoken to her from the Lord.\"",
                    "ref": "Luke 1:45"
                },
                {
                    "text": "And Mary said, \"My soul magnifies the Lord,",
                    "ref": "Luke 1:46"
                },
                {
                    "text": "and my spirit rejoices in God my Savior,",
                    "ref": "Luke 1:47"
                },
                {
                    "text": "for he has regarded the low estate of his handmaiden. For behold, henceforth all generations will call me blessed;\"",
                    "ref": "Luke 1:48"
                }
            ]
        },
        {
            "title": "The Nativity",
            "description": "Jesus is born in a stable in Bethlehem.",
            "virtues": ["Poverty of Spirit", "Detachment"],
            "image": "/images/rosary/annunciation.png",
            "verses": [
                {
                    "text": "And Joseph also went up from Galilee, from the city of Nazareth, to Judea, to the city of David, which is called Bethlehem, because he was of the house and lineage of David,",
                    "ref": "Luke 2:4"
                },
                {
                    "text": "to be enrolled with Mary, his betrothed, who was with child.",
                    "ref": "Luke 2:5"
                },
                {
                    "text": "And while they were there, the time came for her to be delivered.",
                    "ref": "Luke 2:6"
                },
                {
                    "text": "And she gave birth to her first-born son and wrapped him in swaddling cloths, and laid him in a manger, because there was no place for them in the inn.",
                    "ref": "Luke 2:7"
                },
                {
                    "text": "And in that region there were shepherds out in the field, keeping watch over their flock by night.",
                    "ref": "Luke 2:8"
                },
                {
                    "text": "And the angel said to them, \"Be not afraid; for behold, I bring you good news of a great joy which will come to all the people;",
                    "ref": "Luke 2:10"
                },
                {
                    "text": "for to you is born this day in the city of David a Savior, who is Christ the Lord.",
                    "ref": "Luke 2:11"
                },
                {
                    "text": "And this will be a sign for you: you will find a babe wrapped in swaddling cloths and lying in a manger.\"",
                    "ref": "Luke 2:12"
                },
                {
                    "text": "And suddenly there was with the angel a multitude of the heavenly host praising God and saying,",
                    "ref": "Luke 2:13"
                },
                {
                    "text": "\"Glory to God in the highest, and on earth peace among men with whom he is pleased!\"",
                    "ref": "Luke 2:14"
                }
            ]
        },
        {
            "title": "The Presentation",
            "description": "Mary and Joseph present the baby Jesus in the Temple.",
            "virtues": ["Obedience"],
            "image": "/images/rosary/presentation.png",
            "verses": [
                {
                    "text": "And when the time came for their purification according to the law of Moses, they brought him up to Jerusalem to present him to the Lord",
                    "ref": "Luke 2:22"
                },
                {
                    "text": "Now there was a man in Jerusalem, whose name was Simeon, and this man was righteous and devout, looking for the consolation of Israel, and the Holy Spirit was upon him.",
                    "ref": "Luke 2:25"
                },
                {
                    "text": "And it had been revealed to him by the Holy Spirit that he should not see death before he had seen the Lord's Christ.",
                    "ref": "Luke 2:26"
                },
                {
                    "text": "And inspired by the Spirit he came into the temple... he took him up in his arms and blessed God and said,",
                    "ref": "Luke 2:27-28"
                },
                {
                    "text": "\"Lord, now lettest thou thy servant depart in peace, according to thy word;\"",
                    "ref": "Luke 2:29"
                },
                {
                    "text": "\"for mine eyes have seen thy salvation which thou hast prepared in the presence of all peoples,\"",
                    "ref": "Luke 2:30-31"
                },
                {
                    "text": "\"a light for revelation to the Gentiles, and for glory to thy people Israel.\"",
                    "ref": "Luke 2:32"
                },
                {
                    "text": "and Simeon blessed them and said to Mary his mother, \"Behold, this child is set for the fall and rising of many in Israel,\"",
                    "ref": "Luke 2:34"
                },
                {
                    "text": "\"and for a sign that is spoken against (and a sword will pierce through your own soul also), that thoughts out of many hearts may be revealed.\"",
                    "ref": "Luke 2:35"
                },
                {
                    "text": "And coming up at that very hour she [Anna] gave thanks to God, and spoke of him to all who were looking for the redemption of Jerusalem.",
                    "ref": "Luke 2:38"
                }
            ]
        },
        {
            "title": "The Finding in the Temple",
            "description": "Jesus is found teaching the elders in the Temple.",
            "virtues": ["Piety", "Joy of Finding Jesus"],
            "image": "/images/rosary/finding_temple.png",
            "verses": [
                {
                    "text": "Now his parents went to Jerusalem every year at the feast of the Passover.",
                    "ref": "Luke 2:41"
                },
                {
                    "text": "And when he was twelve years old, they went up according to custom;",
                    "ref": "Luke 2:42"
                },
                {
                    "text": "and when the feast was ended, as they were returning, the boy Jesus stayed behind in Jerusalem. His parents did not know it,",
                    "ref": "Luke 2:43"
                },
                {
                    "text": "but supposing him to be in the company they went a day's journey, and they sought him among their kinsfolk and acquaintances;",
                    "ref": "Luke 2:44"
                },
                {
                    "text": "and when they did not find him, they returned to Jerusalem, seeking him.",
                    "ref": "Luke 2:45"
                },
                {
                    "text": "After three days they found him in the temple, sitting among the teachers, listening to them and asking them questions;",
                    "ref": "Luke 2:46"
                },
                {
                    "text": "and all who heard him were amazed at his understanding and his answers.",
                    "ref": "Luke 2:47"
                },
                {
                    "text": "And when they saw him they were astonished; and his mother said to him, \"Son, why have you treated us so? Behold, your father and I have been looking for you anxiously.\"",
                    "ref": "Luke 2:48"
                },
                {
                    "text": "And he said to them, \"How is it that you sought me? Did you not know that I must be in my Father's house?\"",
                    "ref": "Luke 2:49"
                },
                {
                    "text": "And he went down with them and came to Nazareth, and was obedient to them; and his mother kept all these things in her heart.",
                    "ref": "Luke 2:51"
                }
            ]
        }
    ],
    "Luminous": [
        {
            "title": "The Baptism of Jesus",
            "description": "Jesus is baptized by John in the Jordan, and the Holy Spirit descends upon Him.",
            "virtues": ["Fidelity to Baptismal Vows", "Openness to the Holy Spirit"],
            "image": "/images/rosary/annunciation.png",
            "verses": [
                {
                    "text": "John the baptizer appeared in the wilderness, preaching a baptism of repentance for the forgiveness of sins.",
                    "ref": "Mark 1:4"
                },
                {
                    "text": "I have baptized you with water; but he will baptize you with the Holy Spirit.",
                    "ref": "Mark 1:8"
                },
                {
                    "text": "Then Jesus came from Galilee to the Jordan to John, to be baptized by him.",
                    "ref": "Matthew 3:13"
                },
                {
                    "text": "John would have prevented him, saying, \"I need to be baptized by you, and do you come to me?\"",
                    "ref": "Matthew 3:14"
                },
                {
                    "text": "But Jesus answered him, \"Let it be so now; for thus it is fitting for us to fulfil all righteousness.\" Then he consented.",
                    "ref": "Matthew 3:15"
                },
                {
                    "text": "And when Jesus was baptized, he went up immediately from the water, and behold, the heavens were opened",
                    "ref": "Matthew 3:16a"
                },
                {
                    "text": "and he saw the Spirit of God descending like a dove, and alighting on him;",
                    "ref": "Matthew 3:16b"
                },
                {
                    "text": "and lo, a voice from heaven, saying, \"This is my beloved Son, with whom I am well pleased.\"",
                    "ref": "Matthew 3:17"
                },
                {
                    "text": "The next day he saw Jesus coming toward him, and said, \"Behold, the Lamb of God, who takes away the sin of the world!\"",
                    "ref": "John 1:29"
                },
                {
                    "text": "Then Jesus was led up by the Spirit into the wilderness to be tempted by the devil.",
                    "ref": "Matthew 4:1"
                }
            ]
        },
        {
            "title": "The Wedding Feast at Cana",
            "description": "At Mary's request, Jesus performs His first miracle, turning water into wine.",
            "virtues": ["To Jesus through Mary", "Trust in Mary's Intercession"],
            "image": "/images/rosary/wedding_cana.png",
            "verses": [
                {
                    "text": "On the third day there was a marriage at Cana in Galilee, and the mother of Jesus was there;",
                    "ref": "John 2:1"
                },
                {
                    "text": "Jesus also was invited to the marriage, with his disciples.",
                    "ref": "John 2:2"
                },
                {
                    "text": "When the wine failed, the mother of Jesus said to him, \"They have no wine.\"",
                    "ref": "John 2:3"
                },
                {
                    "text": "And Jesus said to her, \"O woman, what have you to do with me? My hour has not yet come.\"",
                    "ref": "John 2:4"
                },
                {
                    "text": "His mother said to the servants, \"Do whatever he tells you.\"",
                    "ref": "John 2:5"
                },
                {
                    "text": "Now six stone jars were standing there, for the Jewish rites of purification, each holding twenty or thirty gallons.",
                    "ref": "John 2:6"
                },
                {
                    "text": "Jesus said to them, \"Fill the jars with water.\" And they filled them up to the brim.",
                    "ref": "John 2:7"
                },
                {
                    "text": "He said to them, \"Now draw some out, and take it to the steward of the feast.\" So they took it.",
                    "ref": "John 2:8"
                },
                {
                    "text": "When the steward of the feast tasted the water now become wine, and did not know where it came from...",
                    "ref": "John 2:9"
                },
                {
                    "text": "This, the first of his signs, Jesus did at Cana in Galilee, and manifested his glory; and his disciples believed in him.",
                    "ref": "John 2:11"
                }
            ]
        },
        {
            "title": "Proclamation of the Kingdom",
            "description": "Jesus announces the coming of the Kingdom of God and calls all to conversion.",
            "virtues": ["Repentance", "Conversion of Heart"],
            "image": "/images/rosary/transfiguration.png",
            "verses": [
                {
                    "text": "Now after John was arrested, Jesus came into Galilee, preaching the gospel of God,",
                    "ref": "Mark 1:14"
                },
                {
                    "text": "and saying, \"The time is fulfilled, and the kingdom of God is at hand; repent, and believe in the gospel.\"",
                    "ref": "Mark 1:15"
                },
                {
                    "text": "\"Blessed are the poor in spirit, for theirs is the kingdom of heaven.\"",
                    "ref": "Matthew 5:3"
                },
                {
                    "text": "\"Blessed are those who mourn, for they shall be comforted.\"",
                    "ref": "Matthew 5:4"
                },
                {
                    "text": "\"Blessed are the meek, for they shall inherit the earth.\"",
                    "ref": "Matthew 5:5"
                },
                {
                    "text": "\"Blessed are those who hunger and thirst for righteousness, for they shall be satisfied.\"",
                    "ref": "Matthew 5:6"
                },
                {
                    "text": "\"Blessed are the merciful, for they shall obtain mercy.\"",
                    "ref": "Matthew 5:7"
                },
                {
                    "text": "\"Blessed are the pure in heart, for they shall see God.\"",
                    "ref": "Matthew 5:8"
                },
                {
                    "text": "And he answered them, \"Go and tell John what you have seen and heard: the blind receive their sight, the lame walk, lepers are cleansed, and the deaf hear, the dead are raised up, the poor have good news preached to them.\"",
                    "ref": "Luke 7:22"
                },
                {
                    "text": "\"And preach as you go, saying, 'The kingdom of heaven is at hand.'\"",
                    "ref": "Matthew 10:7"
                }
            ]
        },
        {
            "title": "The Transfiguration",
            "description": "Jesus is transfigured on Mount Tabor, revealing His divine glory.",
            "virtues": ["Desire for Holiness", "Spiritual Courage"],
            "image": "/images/rosary/annunciation.png",
            "verses": [
                {
                    "text": "And after six days Jesus took with him Peter and James and John his brother, and led them up a high mountain apart.",
                    "ref": "Matthew 17:1"
                },
                {
                    "text": "And he was transfigured before them, and his face shone like the sun, and his garments became white as light.",
                    "ref": "Matthew 17:2"
                },
                {
                    "text": "And behold, there appeared to them Moses and Elijah, talking with him.",
                    "ref": "Matthew 17:3"
                },
                {
                    "text": "And Peter said to Jesus, \"Lord, it is well that we are here; if you wish, I will make three booths here...\"",
                    "ref": "Matthew 17:4"
                },
                {
                    "text": "He was still speaking, when lo, a bright cloud overshadowed them, and a voice from the cloud said,",
                    "ref": "Matthew 17:5a"
                },
                {
                    "text": "\"This is my beloved Son, with whom I am well pleased; listen to him.\"",
                    "ref": "Matthew 17:5b"
                },
                {
                    "text": "When the disciples heard this, they fell on their faces, and were filled with awe.",
                    "ref": "Matthew 17:6"
                },
                {
                    "text": "But Jesus came and touched them, saying, \"Rise, and have no fear.\"",
                    "ref": "Matthew 17:7"
                },
                {
                    "text": "And when they lifted up their eyes, they saw no one but Jesus only.",
                    "ref": "Matthew 17:8"
                },
                {
                    "text": "And as they were coming down the mountain, Jesus commanded them, \"Tell no one the vision, until the Son of man is raised from the dead.\"",
                    "ref": "Matthew 17:9"
                }
            ]
        },
        {
            "title": "Institution of the Eucharist",
            "description": "Jesus institutes the Holy Eucharist at the Last Supper on Holy Thursday.",
            "virtues": ["Eucharistic Adoration", "Active Participation in Mass"],
            "image": "/images/rosary/institution_eucharist.png",
            "verses": [
                {
                    "text": "And he said to them, \"I have earnestly desired to eat this passover with you before I suffer;\"",
                    "ref": "Luke 22:15"
                },
                {
                    "text": "rose from supper, laid aside his garments, and girded himself with a towel... Then he poured water into a basin, and began to wash the disciples' feet.",
                    "ref": "John 13:4-5"
                },
                {
                    "text": "For I received from the Lord what I also delivered to you, that the Lord Jesus on the night when he was betrayed took bread,",
                    "ref": "1 Corinthians 11:23"
                },
                {
                    "text": "Now as they were eating, Jesus took bread, and blessed, and broke it, and gave it to the disciples and said, \"Take, eat; this is my body.\"",
                    "ref": "Matthew 26:26"
                },
                {
                    "text": "And he took a cup, and when he had given thanks he gave it to them, saying, \"Drink of it, all of you;",
                    "ref": "Matthew 26:27"
                },
                {
                    "text": "for this is my blood of the covenant, which is poured out for many for the forgiveness of sins.\"",
                    "ref": "Matthew 26:28"
                },
                {
                    "text": "In the same way also the cup, after supper, saying, \"This cup is the new covenant in my blood. Do this, as often as you drink it, in remembrance of me.\"",
                    "ref": "1 Corinthians 11:25"
                },
                {
                    "text": "And he took bread, and when he had given thanks he broke it and gave it to them, saying, \"This is my body which is given for you. Do this in remembrance of me.\"",
                    "ref": "Luke 22:19"
                },
                {
                    "text": "\"I am the living bread which came down from heaven; if any one eats of this bread, he will live for ever;\"",
                    "ref": "John 6:51a"
                },
                {
                    "text": "\"and the bread which I shall give for the life of the world is my flesh.\"",
                    "ref": "John 6:51b"
                }
            ]
        }
    ],
    "Sorrowful": [
        {
            "title": "Agony in the Garden",
            "description": "Jesus prays in the Garden of Gethsemane on the eve of His passion.",
            "virtues": ["Sorrow for Sin", "Conformity to the Will of God"],
            "image": "/images/rosary/agony_garden.png",
            "verses": [
                {
                    "text": "Then Jesus went with them to a place called Gethsemane, and he said to his disciples, \"Sit here, while I go yonder and pray.\"",
                    "ref": "Matthew 26:36"
                },
                {
                    "text": "And taking with him Peter and the two sons of Zebedee, he began to be sorrowful and troubled.",
                    "ref": "Matthew 26:37"
                },
                {
                    "text": "Then he said to them, \"My soul is very sorrowful, even to death; remain here, and watch with me.\"",
                    "ref": "Matthew 26:38"
                },
                {
                    "text": "And going a little farther he fell on his face and prayed, \"My Father, if it be possible, let this cup pass from me;",
                    "ref": "Matthew 26:39a"
                },
                {
                    "text": "nevertheless, not as I will, but as thou wilt.\"",
                    "ref": "Matthew 26:39b"
                },
                {
                    "text": "And there appeared to him an angel from heaven, strengthening him.",
                    "ref": "Luke 22:43"
                },
                {
                    "text": "And being in an agony he prayed more earnestly; and his sweat became like great drops of blood falling down upon the ground.",
                    "ref": "Luke 22:44"
                },
                {
                    "text": "And he came to the disciples and found them sleeping; and he said to Peter, \"So, could you not watch with me one hour?",
                    "ref": "Matthew 26:40"
                },
                {
                    "text": "Watch and pray that you may not enter into temptation; the spirit indeed is willing, but the flesh is weak.\"",
                    "ref": "Matthew 26:41"
                },
                {
                    "text": "Then he came to the disciples and said to them, \"Are you still sleeping and taking your rest? Behold, the hour is at hand...\"",
                    "ref": "Matthew 26:45"
                }
            ]
        },
        {
            "title": "Scourging at the Pillar",
            "description": "Jesus is cruelly scourged by the Roman soldiers.",
            "virtues": ["Purity", "Mortification of the Senses"],
            "image": "/images/rosary/scourging_pillar.png",
            "verses": [
                {
                    "text": "Pilate said to them, \"Then what shall I do with Jesus who is called Christ?\" They all said, \"Let him be crucified.\"",
                    "ref": "Matthew 27:22"
                },
                {
                    "text": "So when Pilate saw that he was gaining nothing... he took water and washed his hands before the crowd, saying,",
                    "ref": "Matthew 27:24a"
                },
                {
                    "text": "\"I am innocent of this man's blood; see to it yourselves.\"",
                    "ref": "Matthew 27:24b"
                },
                {
                    "text": "And all the people answered, \"His blood be on us and on our children!\"",
                    "ref": "Matthew 27:25"
                },
                {
                    "text": "Then he released for them Barabbas, and having scourged Jesus, delivered him to be crucified.",
                    "ref": "Matthew 27:26"
                },
                {
                    "text": "Then the soldiers of the governor took Jesus into the praetorium, and they gathered the whole battalion before him.",
                    "ref": "Matthew 27:27"
                },
                {
                    "text": "But he was wounded for our transgressions, he was bruised for our iniquities;",
                    "ref": "Isaiah 53:5a"
                },
                {
                    "text": "upon him was the chastisement that made us whole, and with his stripes we are healed.",
                    "ref": "Isaiah 53:5b"
                },
                {
                    "text": "I gave my back to the smiters, and my cheeks to those who pulled out the beard;",
                    "ref": "Isaiah 50:6a"
                },
                {
                    "text": "I hid not my face from shame and spitting.",
                    "ref": "Isaiah 50:6b"
                }
            ]
        },
        {
            "title": "Crowning with Thorns",
            "description": "Jesus is mocked and crowned with a crown of thorns.",
            "virtues": ["Moral Courage", "Reign of Christ in our Hearts"],
            "image": "/images/rosary/crowning_thorns.png",
            "verses": [
                {
                    "text": "And they stripped him and put a scarlet robe upon him,",
                    "ref": "Matthew 27:28"
                },
                {
                    "text": "and plaiting a crown of thorns they put it on his head, and put a reed in his right hand.",
                    "ref": "Matthew 27:29a"
                },
                {
                    "text": "And kneeling before him they mocked him, saying, \"Hail, King of the Jews!\"",
                    "ref": "Matthew 27:29b"
                },
                {
                    "text": "And they spat upon him, and took the reed and struck him on the head.",
                    "ref": "Matthew 27:30"
                },
                {
                    "text": "And they clothed him in a purple cloak, and plaiting a crown of thorns they put it on him.",
                    "ref": "Mark 15:17"
                },
                {
                    "text": "And they began to salute him, \"Hail, King of the Jews!\"",
                    "ref": "Mark 15:18"
                },
                {
                    "text": "And they struck his head with a reed, and spat upon him, and they knelt down in homage to him.",
                    "ref": "Mark 15:19"
                },
                {
                    "text": "So Jesus came out, wearing the crown of thorns and the purple robe. Pilate said to them, \"Behold the man!\"",
                    "ref": "John 19:5"
                },
                {
                    "text": "And when they had mocked him, they stripped him of the robe, and put his own clothes on him, and led him away to crucify him.",
                    "ref": "Matthew 27:31"
                },
                {
                    "text": "And as they led him away, they seized one Simon of Cyrene...",
                    "ref": "Luke 23:26"
                }
            ]
        },
        {
            "title": "Carrying of the Cross",
            "description": "Jesus carries His heavy cross to Calvary.",
            "virtues": ["Patience in Trials", "Grace to carry our Daily Crosses"],
            "image": "/images/rosary/carrying_cross.png",
            "verses": [
                {
                    "text": "So they took Jesus, and he went out, bearing his own cross, to the place called the place of a skull, which is called in Hebrew Golgotha.",
                    "ref": "John 19:17"
                },
                {
                    "text": "And as they led him away, they seized one Simon of Cyrene, who was coming in from the country,",
                    "ref": "Luke 23:26a"
                },
                {
                    "text": "and laid on him the cross, to carry it behind Jesus.",
                    "ref": "Luke 23:26b"
                },
                {
                    "text": "And there followed him a great multitude of the people, and of women who bewailed and lamented him.",
                    "ref": "Luke 23:27"
                },
                {
                    "text": "But Jesus turning to them said, \"Daughters of Jerusalem, do not weep for me, but weep for yourselves and for your children.",
                    "ref": "Luke 23:28"
                },
                {
                    "text": "For if they do this when the wood is green, what will happen when it is dry?\"",
                    "ref": "Luke 23:31"
                },
                {
                    "text": "Two others also, who were criminals, were led away to be put to death with him.",
                    "ref": "Luke 23:32"
                },
                {
                    "text": "And he said to all, \"If any man would come after me, let him deny himself and take up his cross daily and follow me.\"",
                    "ref": "Luke 9:23"
                },
                {
                    "text": "Then Jesus told his disciples, \"If any man would come after me, let him deny himself and take up his cross and follow me.\"",
                    "ref": "Matthew 16:24"
                },
                {
                    "text": "And when they came to the place which is called The Skull, there they crucified him...",
                    "ref": "Luke 23:33"
                }
            ]
        },
        {
            "title": "The Crucifixion",
            "description": "Jesus is nailed to the cross and dies for our salvation.",
            "virtues": ["Salvation", "Self-Sacrifice for Love of God"],
            "image": "/images/rosary/annunciation.png",
            "verses": [
                {
                    "text": "And when they came to the place which is called The Skull, there they crucified him, and the criminals, one on the right and one on the left.",
                    "ref": "Luke 23:33"
                },
                {
                    "text": "And Jesus said, \"Father, forgive them; for they know not what they do.\"",
                    "ref": "Luke 23:34"
                },
                {
                    "text": "And he said, \"Jesus, remember me when you come into your kingdom.\"",
                    "ref": "Luke 23:42"
                },
                {
                    "text": "And he said to him, \"Truly, I say to you, today you will be with me in Paradise.\"",
                    "ref": "Luke 23:43"
                },
                {
                    "text": "When Jesus saw his mother, and the disciple whom he loved standing near, he said to his mother, \"Woman, behold, your son!\"",
                    "ref": "John 19:26"
                },
                {
                    "text": "Then he said to the disciple, \"Behold, your mother!\" And from that hour the disciple took her to his own home.",
                    "ref": "John 19:27"
                },
                {
                    "text": "And when the sixth hour had come, there was darkness over the whole land until the ninth hour.",
                    "ref": "Mark 15:33"
                },
                {
                    "text": "And at the ninth hour Jesus cried with a loud voice, \"Elo-i, Elo-i, la'ma sabach-tha'ni?\" which means, \"My God, my God, why hast thou forsaken me?\"",
                    "ref": "Mark 15:34"
                },
                {
                    "text": "Then Jesus, crying with a loud voice, said, \"Father, into thy hands I commit my spirit!\" And having said this he breathed his last.",
                    "ref": "Luke 23:46"
                },
                {
                    "text": "Now when the centurion saw what had taken place, he praised God, and said, \"Certainly this man was innocent!\"",
                    "ref": "Luke 23:47"
                }
            ]
        }
    ],
    "Glorious": [
        {
            "title": "The Resurrection",
            "description": "Jesus rises from the dead on Easter Sunday.",
            "virtues": ["Faith", "Hope"],
            "image": "/images/rosary/resurrection.png",
            "verses": [
                {
                    "text": "But the angel said to the women, \"Do not be afraid; for I know that you seek Jesus who was crucified.",
                    "ref": "Matthew 28:5"
                },
                {
                    "text": "He is not here; for he has risen, as he said. Come, see the place where he lay.",
                    "ref": "Matthew 28:6"
                },
                {
                    "text": "Then go quickly and tell his disciples that he has risen from the dead... lo, I have told you.\"",
                    "ref": "Matthew 28:7"
                },
                {
                    "text": "On the evening of that day... Jesus came and stood among them",
                    "ref": "John 20:19a"
                },
                {
                    "text": "and said to them, \"Peace be with you.\"",
                    "ref": "John 20:19b"
                },
                {
                    "text": "When he had said this, he showed them his hands and his side. Then the disciples were glad when they saw the Lord.",
                    "ref": "John 20:20"
                },
                {
                    "text": "Jesus said to them again, \"Peace be with you. As the Father has sent me, even so I send you.\"",
                    "ref": "John 20:21"
                },
                {
                    "text": "For I delivered to you as of first importance what I also received, that Christ died for our sins in accordance with the scriptures,",
                    "ref": "1 Corinthians 15:3"
                },
                {
                    "text": "that he was buried, that he was raised on the third day in accordance with the scriptures,",
                    "ref": "1 Corinthians 15:4"
                },
                {
                    "text": "But in fact Christ has been raised from the dead, the first fruits of those who have fallen asleep.",
                    "ref": "1 Corinthians 15:20"
                }
            ]
        },
        {
            "title": "The Ascension",
            "description": "Jesus ascends body and soul into Heaven forty days after His resurrection.",
            "virtues": ["Desire for Heaven", "Hope"],
            "image": "/images/rosary/ascension.png",
            "verses": [
                {
                    "text": "So when they had come together, they asked him, \"Lord, will you at this time restore the kingdom to Israel?\"",
                    "ref": "Acts 1:6"
                },
                {
                    "text": "He said to them, \"It is not for you to know times or seasons which the Father has fixed by his own authority.",
                    "ref": "Acts 1:7"
                },
                {
                    "text": "But you shall receive power when the Holy Spirit has come upon you;",
                    "ref": "Acts 1:8a"
                },
                {
                    "text": "and you shall be my witnesses in Jerusalem and in all Judea and Samaria and to the end of the earth.\"",
                    "ref": "Acts 1:8b"
                },
                {
                    "text": "And when he had said this, as they were looking on, he was lifted up, and a cloud took him out of their sight.",
                    "ref": "Acts 1:9"
                },
                {
                    "text": "And while they were gazing into heaven as he went, behold, two men stood by them in white robes,",
                    "ref": "Acts 1:10"
                },
                {
                    "text": "and said, \"Men of Galilee, why do you stand looking into heaven?",
                    "ref": "Acts 1:11a"
                },
                {
                    "text": "This Jesus, who was taken up from you into heaven, will come in the same way as you saw him go into heaven.\"",
                    "ref": "Acts 1:11b"
                },
                {
                    "text": "So then the Lord Jesus, after he had spoken to them, was taken up into heaven, and sat down at the right hand of God.",
                    "ref": "Mark 16:19"
                },
                {
                    "text": "And they went forth and preached everywhere, while the Lord worked with them...",
                    "ref": "Mark 16:20"
                }
            ]
        },
        {
            "title": "Descent of the Holy Spirit",
            "description": "The Holy Spirit descends as tongues of fire upon Mary and the Apostles on Pentecost.",
            "virtues": ["Love of God", "Holy Wisdom"],
            "image": "/images/rosary/descent_holy_spirit.png",
            "verses": [
                {
                    "text": "When the day of Pentecost had come, they were all together in one place.",
                    "ref": "Acts 2:1"
                },
                {
                    "text": "And suddenly a sound came from heaven like the rush of a mighty wind, and it filled all the house where they were sitting.",
                    "ref": "Acts 2:2"
                },
                {
                    "text": "And there appeared to them tongues as of fire, distributed and resting on each one of them.",
                    "ref": "Acts 2:3"
                },
                {
                    "text": "And they were all filled with the Holy Spirit and began to speak in other tongues, as the Spirit gave them utterance.",
                    "ref": "Acts 2:4"
                },
                {
                    "text": "\"And in the last days it shall be, God declares, that I will pour out my Spirit upon all flesh,\"",
                    "ref": "Acts 2:17a"
                },
                {
                    "text": "\"and your sons and your daughters shall prophesy, and your young men shall see visions, and your old men shall dream dreams;\"",
                    "ref": "Acts 2:17b"
                },
                {
                    "text": "\"And it shall be that whoever calls on the name of the Lord shall be saved.\"",
                    "ref": "Acts 2:21"
                },
                {
                    "text": "But the fruit of the Spirit is love, joy, peace, patience, kindness, goodness, faithfulness,",
                    "ref": "Galatians 5:22"
                },
                {
                    "text": "gentleness, self-control; against such there is no law.",
                    "ref": "Galatians 5:23"
                },
                {
                    "text": "If we live by the Spirit, let us also walk by the Spirit.",
                    "ref": "Galatians 5:25"
                }
            ]
        },
        {
            "title": "Assumption of Mary",
            "description": "Mary is taken body and soul into Heaven at the end of her earthly life.",
            "virtues": ["Devotion to Mary", "Grace of a Happy Death"],
            "image": "/images/rosary/assumption.png",
            "verses": [
                {
                    "text": "My beloved speaks and says to me: \"Arise, my love, my fair one, and come away;\"",
                    "ref": "Song of Solomon 2:10"
                },
                {
                    "text": "\"The fig tree puts forth its figs, and the vines are in blossom; they give forth fragrance. Arise, my love, my fair one, and come away.\"",
                    "ref": "Song of Solomon 2:13"
                },
                {
                    "text": "Then God's temple in heaven was opened, and the ark of his covenant was seen within his temple;",
                    "ref": "Revelation 11:19"
                },
                {
                    "text": "And a great portent appeared in heaven, a woman clothed with the sun, with the moon under her feet, and on her head a crown of twelve stars;",
                    "ref": "Revelation 12:1"
                },
                {
                    "text": "\"for he has regarded the low estate of his handmaiden. For behold, henceforth all generations will call me blessed;\"",
                    "ref": "Luke 1:48"
                },
                {
                    "text": "\"for he who is mighty has done great things for me, and holy is his name.\"",
                    "ref": "Luke 1:49"
                },
                {
                    "text": "Hear, O daughter, consider, and incline your ear; forget your people and your father's house;",
                    "ref": "Psalm 45:10"
                },
                {
                    "text": "and the king will desire your beauty. Since he is your lord, bow to him;",
                    "ref": "Psalm 45:11"
                },
                {
                    "text": "The princess is decked in her chamber with gold-woven robes;",
                    "ref": "Psalm 45:13"
                },
                {
                    "text": "in many-colored robes she is led to the king, with her virgin companions, her escort, in her train.",
                    "ref": "Psalm 45:14"
                }
            ]
        },
        {
            "title": "Coronation of Mary",
            "description": "Mary is crowned Queen of Heaven and Earth by her Son, Jesus.",
            "virtues": ["Trust in Mary's Intercession", "Eternal Happiness"],
            "image": "/images/rosary/coronation.png",
            "verses": [
                {
                    "text": "And a great portent appeared in heaven, a woman clothed with the sun, with the moon under her feet, and on her head a crown of twelve stars;",
                    "ref": "Revelation 12:1"
                },
                {
                    "text": "she was with child and she cried out in her pangs of birth, in anguish for delivery.",
                    "ref": "Revelation 12:2"
                },
                {
                    "text": "And another portent appeared in heaven; behold, a great red dragon, with seven heads and ten horns, and seven diadems upon his heads.",
                    "ref": "Revelation 12:3"
                },
                {
                    "text": "His tail swept down a third of the stars of heaven, and cast them to the earth.",
                    "ref": "Revelation 12:4"
                },
                {
                    "text": "she brought forth a male child, one who is to rule all the nations with a rod of iron, but her child was caught up to God and to his throne,",
                    "ref": "Revelation 12:5"
                },
                {
                    "text": "And Uzziah said to her, \"O daughter, you are blessed by the Most High God above all women on earth;\"",
                    "ref": "Judith 13:18"
                },
                {
                    "text": "\"Your praise will never depart from the hearts of men, who will remember the power of God for ever.\"",
                    "ref": "Judith 13:19"
                },
                {
                    "text": "And when they met her they all blessed her with one accord and said to her, \"You are the exaltation of Jerusalem, you are the great glory of Israel, you are the great pride of our nation!\"",
                    "ref": "Judith 15:9"
                },
                {
                    "text": "And I heard a loud voice in heaven, saying, \"Now the salvation and the power and the kingdom of our God and the authority of his Christ have come...\"",
                    "ref": "Revelation 12:10"
                },
                {
                    "text": "\"And they have conquered him by the blood of the Lamb and by the word of their testimony, for they loved not their lives even unto death.\"",
                    "ref": "Revelation 12:11"
                }
            ]
        }
    ]
};

export const generateRosarySteps = (
    mysterySet: string,
    beforePrayerIds: string[] = [],
    afterPrayerIds: string[] = []
): RosaryStep[] => {
    const steps: RosaryStep[] = [];
    const mysteries = ROSARY_DATA[mysterySet] || [];

    // 1. Intro
    steps.push({ type: 'opening', prayerId: 'sign-of-the-cross' });

    // Prepend custom prayers
    beforePrayerIds.forEach(prayerId => {
        steps.push({ type: 'intro', prayerId });
    });

    steps.push({ type: 'intro', prayerId: 'apostles-creed' });
    steps.push({ type: 'intro', prayerId: 'our-father' });
    steps.push({ type: 'intro', prayerId: 'hail-mary', beadNumber: 1 });
    steps.push({ type: 'intro', prayerId: 'hail-mary', beadNumber: 2 });
    steps.push({ type: 'intro', prayerId: 'hail-mary', beadNumber: 3 });
    steps.push({ type: 'intro', prayerId: 'glory-be' });

    // 2. Mysteries
    mysteries.forEach((mystery, index) => {
        const decadeNum = index + 1;

        // Mystery Announcement
        steps.push({
            type: 'mystery-header',
            prayerId: '',
            mysteryTitle: mystery.title,
            decadeNumber: decadeNum,
            description: mystery.description,
            virtues: mystery.virtues
        });

        // Our Father
        steps.push({ type: 'decade-start', prayerId: 'our-father', decadeNumber: decadeNum });

        // 10 Hail Marys
        for (let i = 0; i < 10; i++) {
            steps.push({
                type: 'decade-bead',
                prayerId: 'hail-mary',
                beadNumber: i + 1,
                decadeNumber: decadeNum,
                verse: mystery.verses[i]
            });
        }

        // Glory Be + Fatima
        steps.push({ type: 'decade-end', prayerId: 'glory-be', decadeNumber: decadeNum });
        steps.push({ type: 'decade-end', prayerId: 'fatima-prayer', decadeNumber: decadeNum });
    });

    // 3. Closing
    steps.push({ type: 'closing', prayerId: 'hail-holy-queen' });
    steps.push({ type: 'closing', prayerId: 'rosary-closing-prayer' });

    // Append custom prayers
    afterPrayerIds.forEach(prayerId => {
        steps.push({ type: 'closing', prayerId });
    });

    steps.push({ type: 'opening', prayerId: 'sign-of-the-cross' });

    // 4. Completion Marker
    steps.push({
        type: 'closing',
        prayerId: 'completion-marker',
        title: 'Rosary Completed',
        content: 'You have completed the Rosary devotion. May the peace of Christ rule in your heart.',
        latin: 'Rosarium finitum est. Pax Christi abundet in corde tuo.'
    });

    return steps;
};
