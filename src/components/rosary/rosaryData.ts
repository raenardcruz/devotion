export type Verse = {
    text: string;
    ref: string;
};

export type Mystery = {
    title: string;
    focus: string;
    verses: Verse[];
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
};

export const ROSARY_DATA: RosaryData = {
    "Joyful": [
        {
            "title": "The Annunciation",
            "focus": "The Angel Gabriel asks Mary to be the Mother of God.",
            "verses": [
                {
                    "text": "In the sixth month the angel Gabriel was sent by God to a town in Galilee called Nazareth,",
                    "ref": "Luke 1:26"
                },
                {
                    "text": "to a virgin engaged to a man whose name was Joseph, of the house of David. The virgin's name was Mary.",
                    "ref": "Luke 1:27"
                },
                {
                    "text": "And he came to her and said, \"Greetings, favored one! The Lord is with you.\"",
                    "ref": "Luke 1:28"
                },
                {
                    "text": "But she was much perplexed by his words and pondered what sort of greeting this might be.",
                    "ref": "Luke 1:29"
                },
                {
                    "text": "The angel said to her, \"Do not be afraid, Mary, for you have found favor with God.\"",
                    "ref": "Luke 1:30"
                },
                {
                    "text": "And now, you will conceive in your womb and bear a son, and you will name him Jesus.",
                    "ref": "Luke 1:31"
                },
                {
                    "text": "He will be great, and will be called the Son of the Most High, and the Lord God will give to him the throne of his ancestor David.",
                    "ref": "Luke 1:32"
                },
                {
                    "text": "Mary said to the angel, \"How can this be, since I am a virgin?\"",
                    "ref": "Luke 1:34"
                },
                {
                    "text": "The angel said to her, \"The Holy Spirit will come upon you, and the power of the Most High will overshadow you; therefore the child to be born will be holy; he will be called Son of God.\"",
                    "ref": "Luke 1:35"
                },
                {
                    "text": "Then Mary said, \"Here am I, the servant of the Lord; let it be with me according to your word.\" Then the angel departed from her.",
                    "ref": "Luke 1:38"
                }
            ]
        },
        {
            "title": "The Visitation",
            "focus": "Mary visits her cousin Elizabeth.",
            "verses": [
                {
                    "text": "In those days Mary set out and went with haste to a Judean town in the hill country,",
                    "ref": "Luke 1:39"
                },
                {
                    "text": "where she entered the house of Zechariah and greeted Elizabeth.",
                    "ref": "Luke 1:40"
                },
                {
                    "text": "When Elizabeth heard Mary's greeting, the child leaped in her womb. And Elizabeth was filled with the Holy Spirit",
                    "ref": "Luke 1:41"
                },
                {
                    "text": "and exclaimed with a loud cry, \"Blessed are you among women, and blessed is the fruit of your womb.\"",
                    "ref": "Luke 1:42"
                },
                {
                    "text": "And why has this happened to me, that the mother of my Lord comes to me?",
                    "ref": "Luke 1:43"
                },
                {
                    "text": "For as soon as I heard the sound of your greeting, the child in my womb leaped for joy.",
                    "ref": "Luke 1:44"
                },
                {
                    "text": "And blessed is she who believed that there would be a fulfillment of what was spoken to her by the Lord.\"",
                    "ref": "Luke 1:45"
                },
                {
                    "text": "And Mary said, \"My soul magnifies the Lord,\"",
                    "ref": "Luke 1:46"
                },
                {
                    "text": "and my spirit rejoices in God my Savior,\"",
                    "ref": "Luke 1:47"
                },
                {
                    "text": "for he has looked with favor on the lowliness of his servant. Surely, from now on all generations will call me blessed;",
                    "ref": "Luke 1:48"
                }
            ]
        },
        {
            "title": "The Nativity",
            "focus": "Jesus is born in a stable in Bethlehem.",
            "verses": [
                {
                    "text": "Joseph also went from the town of Nazareth in Galilee to Judea, to the city of David called Bethlehem, because he was descended from the house and family of David.",
                    "ref": "Luke 2:4"
                },
                {
                    "text": "He went to be registered with Mary, to whom he was engaged and who was expecting a child.",
                    "ref": "Luke 2:5"
                },
                {
                    "text": "While they were there, the time came for her to deliver her child.",
                    "ref": "Luke 2:6"
                },
                {
                    "text": "And she gave birth to her firstborn son and wrapped him in bands of cloth, and laid him in a manger, because there was no place for them in the inn.",
                    "ref": "Luke 2:7"
                },
                {
                    "text": "In that region there were shepherds living in the fields, keeping watch over their flock by night.",
                    "ref": "Luke 2:8"
                },
                {
                    "text": "But the angel said to them, \"Do not be afraid; for see—I am bringing you good news of great joy for all the people:",
                    "ref": "Luke 2:10"
                },
                {
                    "text": "to you is born this day in the city of David a Savior, who is the Messiah, the Lord.\"",
                    "ref": "Luke 2:11"
                },
                {
                    "text": "\"This will be a sign for you: you will find a child wrapped in bands of cloth and lying in a manger.\"",
                    "ref": "Luke 2:12"
                },
                {
                    "text": "And suddenly there was with the angel a multitude of the heavenly host, praising God and saying,",
                    "ref": "Luke 2:13"
                },
                {
                    "text": "\"Glory to God in the highest heaven, and on earth peace among those whom he favors!\"",
                    "ref": "Luke 2:14"
                }
            ]
        },
        {
            "title": "The Presentation",
            "focus": "Mary and Joseph take baby Jesus to the Temple.",
            "verses": [
                {
                    "text": "When the time came for their purification according to the law of Moses, they brought him up to Jerusalem to present him to the Lord",
                    "ref": "Luke 2:22"
                },
                {
                    "text": "Now there was a man in Jerusalem whose name was Simeon; this man was righteous and devout, looking forward to the consolation of Israel, and the Holy Spirit rested on him.",
                    "ref": "Luke 2:25"
                },
                {
                    "text": "It had been revealed to him by the Holy Spirit that he would not see death before he had seen the Lord's Messiah.",
                    "ref": "Luke 2:26"
                },
                {
                    "text": "Guided by the Spirit, Simeon came into the temple... he took him in his arms and praised God, saying,",
                    "ref": "Luke 2:27-28"
                },
                {
                    "text": "\"Master, now you are dismissing your servant in peace, according to your word;\"",
                    "ref": "Luke 2:29"
                },
                {
                    "text": "\"for my eyes have seen your salvation, which you have prepared in the presence of all peoples,\"",
                    "ref": "Luke 2:30-31"
                },
                {
                    "text": "\"a light for revelation to the Gentiles and for glory to your people Israel.\"",
                    "ref": "Luke 2:32"
                },
                {
                    "text": "Then Simeon blessed them and said to his mother Mary, \"This child is destined for the falling and the rising of many in Israel,\"",
                    "ref": "Luke 2:34"
                },
                {
                    "text": "\"and to be a sign that will be opposed so that the inner thoughts of many will be revealed—and a sword will pierce your own soul too.\"",
                    "ref": "Luke 2:35"
                },
                {
                    "text": "At that moment she [Anna] came, and began to praise God and to speak about the child to all who were looking for the redemption of Jerusalem.",
                    "ref": "Luke 2:38"
                }
            ]
        },
        {
            "title": "The Finding in the Temple",
            "focus": "Jesus is found teaching the elders in the Temple.",
            "verses": [
                {
                    "text": "Now every year his parents went to Jerusalem for the festival of the Passover.",
                    "ref": "Luke 2:41"
                },
                {
                    "text": "And when he was twelve years old, they went up as usual for the festival.",
                    "ref": "Luke 2:42"
                },
                {
                    "text": "When the festival was ended and they started to return, the boy Jesus stayed behind in Jerusalem, but his parents did not know it.",
                    "ref": "Luke 2:43"
                },
                {
                    "text": "Assuming that he was in the group of travelers, they went a day's journey. Then they started to look for him among their relatives and friends.",
                    "ref": "Luke 2:44"
                },
                {
                    "text": "When they did not find him, they returned to Jerusalem to search for him.",
                    "ref": "Luke 2:45"
                },
                {
                    "text": "After three days they found him in the temple, sitting among the teachers, listening to them and asking them questions.",
                    "ref": "Luke 2:46"
                },
                {
                    "text": "And all who heard him were amazed at his understanding and his answers.",
                    "ref": "Luke 2:47"
                },
                {
                    "text": "His mother said to him, \"Child, why have you treated us like this? Look, your father and I have been searching for you in great anxiety.\"",
                    "ref": "Luke 2:48"
                },
                {
                    "text": "He said to them, \"Why were you searching for me? Did you not know that I must be in my Father’s house?\"",
                    "ref": "Luke 2:49"
                },
                {
                    "text": "Then he went down with them and came to Nazareth, and was obedient to them. His mother treasured all these things in her heart.",
                    "ref": "Luke 2:51"
                }
            ]
        }
    ],
    "Luminous": [
        {
            "title": "The Baptism of Jesus",
            "focus": "Openness to the Holy Spirit.",
            "verses": [
                {
                    "text": "John the baptizer appeared in the wilderness, proclaiming a baptism of repentance for the forgiveness of sins.",
                    "ref": "Mark 1:4"
                },
                {
                    "text": "I have baptized you with water; but he will baptize you with the Holy Spirit.",
                    "ref": "Mark 1:8"
                },
                {
                    "text": "In those days Jesus came from Nazareth of Galilee and was baptized by John in the Jordan.",
                    "ref": "Mark 1:9"
                },
                {
                    "text": "And just as he was coming up out of the water, he saw the heavens torn apart and the Spirit descending like a dove on him.",
                    "ref": "Mark 1:10"
                },
                {
                    "text": "And a voice came from heaven, \"You are my Son, the Beloved; with you I am well pleased.\"",
                    "ref": "Mark 1:11"
                },
                {
                    "text": "The next day he saw Jesus coming toward him and declared, \"Here is the Lamb of God who takes away the sin of the world!\"",
                    "ref": "John 1:29"
                },
                {
                    "text": "And John testified, \"I saw the Spirit descending from heaven like a dove, and it remained on him.\"",
                    "ref": "John 1:32"
                },
                {
                    "text": "And I myself have seen and have testified that this is the Son of God.\"",
                    "ref": "John 1:34"
                },
                {
                    "text": "Then Jesus was led up by the Spirit into the wilderness to be tempted by the devil.",
                    "ref": "Matthew 4:1"
                },
                {
                    "text": "Go therefore and make disciples of all nations, baptizing them in the name of the Father and of the Son and of the Holy Spirit,",
                    "ref": "Matthew 28:19"
                }
            ]
        },
        {
            "title": "The Wedding Feast at Cana",
            "focus": "To Jesus through Mary.",
            "verses": [
                {
                    "text": "On the third day there was a wedding in Cana of Galilee, and the mother of Jesus was there.",
                    "ref": "John 2:1"
                },
                {
                    "text": "Jesus and his disciples had also been invited to the wedding.",
                    "ref": "John 2:2"
                },
                {
                    "text": "When the wine gave out, the mother of Jesus said to him, \"They have no wine.\"",
                    "ref": "John 2:3"
                },
                {
                    "text": "And Jesus said to her, \"Woman, what concern is that to you and to me? My hour has not yet come.\"",
                    "ref": "John 2:4"
                },
                {
                    "text": "His mother said to the servants, \"Do whatever he tells you.\"",
                    "ref": "John 2:5"
                },
                {
                    "text": "Now standing there were six stone water jars for the Jewish rites of purification, each holding twenty or thirty gallons.",
                    "ref": "John 2:6"
                },
                {
                    "text": "Jesus said to them, \"Fill the jars with water.\" And they filled them up to the brim.",
                    "ref": "John 2:7"
                },
                {
                    "text": "He said to them, \"Now draw some out, and take it to the chief steward.\" So they took it.",
                    "ref": "John 2:8"
                },
                {
                    "text": "When the steward tasted the water that had become wine, and did not know where it came from...",
                    "ref": "John 2:9"
                },
                {
                    "text": "Jesus did this, the first of his signs, in Cana of Galilee, and revealed his glory; and his disciples believed in him.",
                    "ref": "John 2:11"
                }
            ]
        },
        {
            "title": "Proclamation of the Kingdom",
            "focus": "Repentance and trust in the Gospel.",
            "verses": [
                {
                    "text": "\"Blessed are the poor in spirit, for theirs is the kingdom of heaven.\"",
                    "ref": "Matthew 5:3"
                },
                {
                    "text": "\"Blessed are those who mourn, for they will be comforted.\"",
                    "ref": "Matthew 5:4"
                },
                {
                    "text": "\"Blessed are the meek, for they will inherit the earth.\"",
                    "ref": "Matthew 5:5"
                },
                {
                    "text": "\"Blessed are those who hunger and thirst for righteousness, for they will be filled.\"",
                    "ref": "Matthew 5:6"
                },
                {
                    "text": "\"Blessed are the merciful, for they will receive mercy.\"",
                    "ref": "Matthew 5:7"
                },
                {
                    "text": "\"Blessed are the pure in heart, for they will see God.\"",
                    "ref": "Matthew 5:8"
                },
                {
                    "text": "\"Blessed are the peacemakers, for they will be called children of God.\"",
                    "ref": "Matthew 5:9"
                },
                {
                    "text": "\"Blessed are those who are persecuted for righteousness' sake, for theirs is the kingdom of heaven.\"",
                    "ref": "Matthew 5:10"
                },
                {
                    "text": "\"Blessed are you when people revile you and persecute you and utter all kinds of evil against you falsely on my account.\"",
                    "ref": "Matthew 5:11"
                },
                {
                    "text": "\"Rejoice and be glad, for your reward is great in heaven, for in the same way they persecuted the prophets who were before you.\"",
                    "ref": "Matthew 5:12"
                }
            ]
        },
        {
            "title": "The Transfiguration",
            "focus": "Spiritual courage and holiness.",
            "verses": [
                {
                    "text": "Six days later, Jesus took with him Peter and James and his brother John and led them up a high mountain, by themselves.",
                    "ref": "Matthew 17:1"
                },
                {
                    "text": "And he was transfigured before them, and his face shone like the sun, and his clothes became dazzling white.",
                    "ref": "Matthew 17:2"
                },
                {
                    "text": "Suddenly there appeared to them Moses and Elijah, talking with him.",
                    "ref": "Matthew 17:3"
                },
                {
                    "text": "Then Peter said to Jesus, \"Lord, it is good for us to be here; if you wish, I will make three dwellings here...\"",
                    "ref": "Matthew 17:4"
                },
                {
                    "text": "While he was still speaking, suddenly a bright cloud overshadowed them, and from the cloud a voice said,",
                    "ref": "Matthew 17:5"
                },
                {
                    "text": "\"This is my Son, the Beloved; with him I am well pleased; listen to him!\"",
                    "ref": "Matthew 17:5"
                },
                {
                    "text": "When the disciples heard this, they fell to the ground and were overcome by fear.",
                    "ref": "Matthew 17:6"
                },
                {
                    "text": "But Jesus came and touched them, saying, \"Get up and do not be afraid.\"",
                    "ref": "Matthew 17:7"
                },
                {
                    "text": "And when they looked up, they saw no one except Jesus himself alone.",
                    "ref": "Matthew 17:8"
                },
                {
                    "text": "As they were coming down the mountain, Jesus ordered them, \"Tell no one about the vision until after the Son of Man has been raised from the dead.\"",
                    "ref": "Matthew 17:9"
                }
            ]
        },
        {
            "title": "Institution of the Eucharist",
            "focus": "Real presence of Christ.",
            "verses": [
                {
                    "text": "He said to them, \"I have eagerly desired to eat this Passover with you before I suffer;\"",
                    "ref": "Luke 22:15"
                },
                {
                    "text": "got up from the table, took off his outer robe, and tied a towel around himself... and began to wash the disciples' feet.",
                    "ref": "John 13:4-5"
                },
                {
                    "text": "So if I, your Lord and Teacher, have washed your feet, you also ought to wash one another's feet.",
                    "ref": "John 13:14"
                },
                {
                    "text": "While they were eating, Jesus took a loaf of bread, and after blessing it he broke it, gave it to the disciples, and said, \"Take, eat; this is my body.\"",
                    "ref": "Matthew 26:26"
                },
                {
                    "text": "Then he took a cup, and after giving thanks he gave it to them, saying, \"Drink from it, all of you;\"",
                    "ref": "Matthew 26:27"
                },
                {
                    "text": "for this is my blood of the covenant, which is poured out for many for the forgiveness of sins.",
                    "ref": "Matthew 26:28"
                },
                {
                    "text": "Then he took a loaf of bread, and when he had given thanks, he broke it and gave it to them, saying,",
                    "ref": "Luke 22:19"
                },
                {
                    "text": "\"This is my body, which is given for you. Do this in remembrance of me.\"",
                    "ref": "Luke 22:19"
                },
                {
                    "text": "\"I am the living bread that came down from heaven. Whoever eats of this bread will live forever;\"",
                    "ref": "John 6:51"
                },
                {
                    "text": "\"and the bread that I will give for the life of the world is my flesh.\"",
                    "ref": "John 6:51"
                }
            ]
        }
    ],
    "Sorrowful": [
        {
            "title": "Agony in the Garden",
            "focus": "Contrition and conformity to God's will.",
            "verses": [
                {
                    "text": "Then Jesus went with them to a place called Gethsemane; and he said to his disciples, \"Sit here while I go over there and pray.\"",
                    "ref": "Matthew 26:36"
                },
                {
                    "text": "He took with him Peter and the two sons of Zebedee, and began to be grieved and agitated.",
                    "ref": "Matthew 26:37"
                },
                {
                    "text": "Then he said to them, \"I am deeply grieved, even to death; remain here, and stay awake with me.\"",
                    "ref": "Matthew 26:38"
                },
                {
                    "text": "And going a little farther, he threw himself on the ground and prayed, \"My Father, if it is possible, let this cup pass from me;\"",
                    "ref": "Matthew 26:39"
                },
                {
                    "text": "yet not what I want but what you want.\"",
                    "ref": "Matthew 26:39"
                },
                {
                    "text": "Then an angel from heaven appeared to him and gave him strength.",
                    "ref": "Luke 22:43"
                },
                {
                    "text": "In his anguish he prayed more earnestly, and his sweat became like great drops of blood falling down on the ground.",
                    "ref": "Luke 22:44"
                },
                {
                    "text": "Then he came to the disciples and found them sleeping; and he said to Peter, \"So, could you not stay awake with me one hour?\"",
                    "ref": "Matthew 26:40"
                },
                {
                    "text": "Stay awake and pray that you may not come into the time of trial; the spirit indeed is willing, but the flesh is weak.\"",
                    "ref": "Matthew 26:41"
                },
                {
                    "text": "Then he came to the disciples and said to them, \"Are you still sleeping and taking your rest? See, the hour is at hand...\"",
                    "ref": "Matthew 26:45"
                }
            ]
        },
        {
            "title": "Scourging at the Pillar",
            "focus": "Purity and mortification.",
            "verses": [
                {
                    "text": "Pilate said to them, \"Then what should I do with Jesus who is called the Messiah?\" All of them said, \"Let him be crucified!\"",
                    "ref": "Matthew 27:22"
                },
                {
                    "text": "So when Pilate saw that he could do nothing... he took some water and washed his hands before the crowd, saying,",
                    "ref": "Matthew 27:24"
                },
                {
                    "text": "\"I am innocent of this man's blood; see to it yourselves.\"",
                    "ref": "Matthew 27:24"
                },
                {
                    "text": "Then the people as a whole answered, \"His blood be on us and on our children!\"",
                    "ref": "Matthew 27:25"
                },
                {
                    "text": "So he released Barabbas for them; and after flogging Jesus, he handed him over to be crucified.",
                    "ref": "Matthew 27:26"
                },
                {
                    "text": "Then the soldiers of the governor took Jesus into the governor's headquarters, and they gathered the whole cohort around him.",
                    "ref": "Matthew 27:27"
                },
                {
                    "text": "But he was wounded for our transgressions, crushed for our iniquities;",
                    "ref": "Isaiah 53:5"
                },
                {
                    "text": "upon him was the punishment that made us whole, and by his bruises we are healed.",
                    "ref": "Isaiah 53:5"
                },
                {
                    "text": "I gave my back to those who struck me, and my cheeks to those who pulled out the beard;",
                    "ref": "Isaiah 50:6"
                },
                {
                    "text": "I did not hide my face from insult and spitting.",
                    "ref": "Isaiah 50:6"
                }
            ]
        },
        {
            "title": "Crowning with Thorns",
            "focus": "Moral courage.",
            "verses": [
                {
                    "text": "They stripped him and put a scarlet robe on him,",
                    "ref": "Matthew 27:28"
                },
                {
                    "text": "and after twisting some thorns into a crown, they put it on his head. They put a reed in his right hand",
                    "ref": "Matthew 27:29"
                },
                {
                    "text": "and knelt before him and mocked him, saying, \"Hail, King of the Jews!\"",
                    "ref": "Matthew 27:29"
                },
                {
                    "text": "They spat on him, and took the reed and struck him on the head.",
                    "ref": "Matthew 27:30"
                },
                {
                    "text": "So Jesus came out, wearing the crown of thorns and the purple robe. Pilate said to them, \"Here is the man!\"",
                    "ref": "John 19:5"
                },
                {
                    "text": "When the chief priests and the police saw him, they shouted, \"Crucify him! Crucify him!\"",
                    "ref": "John 19:6"
                },
                {
                    "text": "They cried out, \"Away with him! Away with him! Crucify him!\"",
                    "ref": "John 19:15"
                },
                {
                    "text": "Pilate asked them, \"Shall I crucify your King?\" The chief priests answered, \"We have no king but Emperor Caesar.\"",
                    "ref": "John 19:15"
                },
                {
                    "text": "After mocking him, they stripped him of the robe and put his own clothes on him. Then they led him away to crucify him.",
                    "ref": "Matthew 27:31"
                },
                {
                    "text": "As they led him away, they seized a man, Simon of Cyrene...",
                    "ref": "Luke 23:26"
                }
            ]
        },
        {
            "title": "Carrying of the Cross",
            "focus": "Patience in trials.",
            "verses": [
                {
                    "text": "So they took Jesus; and carrying the cross by himself, he went out to what is called The Place of the Skull, which in Hebrew is called Golgotha.",
                    "ref": "John 19:16-17"
                },
                {
                    "text": "As they led him away, they seized a man, Simon of Cyrene, who was coming from the country,",
                    "ref": "Luke 23:26"
                },
                {
                    "text": "and they laid the cross on him, and made him carry it behind Jesus.",
                    "ref": "Luke 23:26"
                },
                {
                    "text": "A great number of the people followed him, and among them were women who were beating their breasts and wailing for him.",
                    "ref": "Luke 23:27"
                },
                {
                    "text": "But Jesus turned to them and said, \"Daughters of Jerusalem, do not weep for me, but weep for yourselves and for your children.\"",
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
                    "text": "Then he said to them all, \"If any want to become my followers, let them deny themselves and take up their cross daily and follow me.\"",
                    "ref": "Luke 9:23"
                },
                {
                    "text": "and whoever does not take up the cross and follow me is not worthy of me.",
                    "ref": "Matthew 10:38"
                },
                {
                    "text": "And when they came to the place that is called The Skull, they crucified Jesus there...",
                    "ref": "Luke 23:33"
                }
            ]
        },
        {
            "title": "The Crucifixion",
            "focus": "Self-sacrifice and forgiveness.",
            "verses": [
                {
                    "text": "When they came to the place that is called The Skull, they crucified Jesus there with the criminals, one on his right and one on his left.",
                    "ref": "Luke 23:33"
                },
                {
                    "text": "Then Jesus said, \"Father, forgive them; for they do not know what they are doing.\"",
                    "ref": "Luke 23:34"
                },
                {
                    "text": "Then he said, \"Jesus, remember me when you come into your kingdom.\"",
                    "ref": "Luke 23:42"
                },
                {
                    "text": "He replied, \"Truly I tell you, today you will be with me in Paradise.\"",
                    "ref": "Luke 23:43"
                },
                {
                    "text": "When Jesus saw his mother and the disciple whom he loved standing beside her, he said to his mother, \"Woman, here is your son.\"",
                    "ref": "John 19:26"
                },
                {
                    "text": "Then he said to the disciple, \"Here is your mother.\"",
                    "ref": "John 19:27"
                },
                {
                    "text": "When it was noon, darkness came over the whole land until three in the afternoon.",
                    "ref": "Mark 15:33"
                },
                {
                    "text": "At three o'clock Jesus cried out with a loud voice, \"Eloi, Eloi, lema sabachthani?\" which means, \"My God, my God, why have you forsaken me?\"",
                    "ref": "Mark 15:34"
                },
                {
                    "text": "Then Jesus, crying with a loud voice, said, \"Father, into your hands I commend my spirit.\" Having said this, he breathed his last.",
                    "ref": "Luke 23:46"
                },
                {
                    "text": "When the centurion saw what had taken place, he praised God and said, \"Certainly this man was innocent.\"",
                    "ref": "Luke 23:47"
                }
            ]
        }
    ],
    "Glorious": [
        {
            "title": "The Resurrection",
            "focus": "Faith and hope.",
            "verses": [
                {
                    "text": "But the angel said to the women, \"Do not be afraid; I know that you are looking for Jesus who was crucified.\"",
                    "ref": "Matthew 28:5"
                },
                {
                    "text": "\"He is not here; for he has been raised, as he said. Come, see the place where he lay.\"",
                    "ref": "Matthew 28:6"
                },
                {
                    "text": "\"Then go quickly and tell his disciples, 'He has been raised from the dead... there you will see him.'\"",
                    "ref": "Matthew 28:7"
                },
                {
                    "text": "When it was evening on that day... Jesus came and stood among them",
                    "ref": "John 20:19"
                },
                {
                    "text": "and said to them, \"Peace be with you.\"",
                    "ref": "John 20:19"
                },
                {
                    "text": "After he said this, he showed them his hands and his side. Then the disciples rejoiced when they saw the Lord.",
                    "ref": "John 20:20"
                },
                {
                    "text": "Jesus said to them again, \"Peace be with you. As the Father has sent me, so I send you.\"",
                    "ref": "John 20:21"
                },
                {
                    "text": "Thomas answered him, \"My Lord and my God!\"",
                    "ref": "John 20:28"
                },
                {
                    "text": "Jesus said to him, \"Have you believed because you have seen me? Blessed are those who have not seen and yet have come to believe.\"",
                    "ref": "John 20:29"
                },
                {
                    "text": "Now Jesus did many other signs in the presence of his disciples...",
                    "ref": "John 20:30"
                }
            ]
        },
        {
            "title": "The Ascension",
            "focus": "Hope and leadership.",
            "verses": [
                {
                    "text": "So when they had come together, they asked him, \"Lord, is this the time when you will restore the kingdom to Israel?\"",
                    "ref": "Acts 1:6"
                },
                {
                    "text": "He replied, \"It is not for you to know the times or periods that the Father has set by his own authority.\"",
                    "ref": "Acts 1:7"
                },
                {
                    "text": "\"But you will receive power when the Holy Spirit has come upon you;\"",
                    "ref": "Acts 1:8"
                },
                {
                    "text": "\"and you will be my witnesses in Jerusalem, in all Judea and Samaria, and to the ends of the earth.\"",
                    "ref": "Acts 1:8"
                },
                {
                    "text": "When he had said this, as they were watching, he was lifted up, and a cloud took him out of their sight.",
                    "ref": "Acts 1:9"
                },
                {
                    "text": "While he was going and they were gazing up toward heaven, suddenly two men in white robes stood by them.",
                    "ref": "Acts 1:10"
                },
                {
                    "text": "They said, \"Men of Galilee, why do you stand looking up toward heaven?\"",
                    "ref": "Acts 1:11"
                },
                {
                    "text": "\"This Jesus, who has been taken up from you into heaven, will come in the same way as you saw him go into heaven.\"",
                    "ref": "Acts 1:11"
                },
                {
                    "text": "So then the Lord Jesus, after he had spoken to them, was taken up into heaven and sat down at the right hand of God.",
                    "ref": "Mark 16:19"
                },
                {
                    "text": "And they went out and proclaimed the good news everywhere, while the Lord worked with them...",
                    "ref": "Mark 16:20"
                }
            ]
        },
        {
            "title": "Descent of the Holy Spirit",
            "focus": "Love of God and holy wisdom.",
            "verses": [
                {
                    "text": "When the day of Pentecost had come, they were all together in one place.",
                    "ref": "Acts 2:1"
                },
                {
                    "text": "And suddenly from heaven there came a sound like the rush of a violent wind, and it filled the entire house where they were sitting.",
                    "ref": "Acts 2:2"
                },
                {
                    "text": "Divided tongues, as of fire, appeared among them, and a tongue rested on each of them.",
                    "ref": "Acts 2:3"
                },
                {
                    "text": "All of them were filled with the Holy Spirit and began to speak in other languages, as the Spirit gave them ability.",
                    "ref": "Acts 2:4"
                },
                {
                    "text": "\"In the last days it will be, God declares, that I will pour out my Spirit upon all flesh,\"",
                    "ref": "Acts 2:17"
                },
                {
                    "text": "\"and your sons and your daughters shall prophesy, and your young men shall see visions, and your old men shall dream dreams.\"",
                    "ref": "Acts 2:17"
                },
                {
                    "text": "\"Then everyone who calls on the name of the Lord shall be saved.\"",
                    "ref": "Acts 2:21"
                },
                {
                    "text": "Create in me a clean heart, O God, and put a new and right spirit within me.",
                    "ref": "Psalm 51:10"
                },
                {
                    "text": "Do not cast me away from your presence, and do not take your holy spirit from me.",
                    "ref": "Psalm 51:11"
                },
                {
                    "text": "Restore to me the joy of your salvation, and sustain in me a willing spirit.",
                    "ref": "Psalm 51:12"
                }
            ]
        },
        {
            "title": "Assumption of Mary",
            "focus": "Devotion to Mary.",
            "verses": [
                {
                    "text": "And Mary said, \"My soul magnifies the Lord,\"",
                    "ref": "Luke 1:46"
                },
                {
                    "text": "\"and my spirit rejoices in God my Savior,\"",
                    "ref": "Luke 1:47"
                },
                {
                    "text": "\"for he has looked with favor on the lowliness of his servant. Surely, from now on all generations will call me blessed;\"",
                    "ref": "Luke 1:48"
                },
                {
                    "text": "\"for the Mighty One has done great things for me, and holy is his name.\"",
                    "ref": "Luke 1:49"
                },
                {
                    "text": "\"His mercy is for those who fear him from generation to generation.\"",
                    "ref": "Luke 1:50"
                },
                {
                    "text": "\"He has shown strength with his arm; he has scattered the proud in the thoughts of their hearts.\"",
                    "ref": "Luke 1:51"
                },
                {
                    "text": "\"He has brought down the powerful from their thrones, and lifted up the lowly;\"",
                    "ref": "Luke 1:52"
                },
                {
                    "text": "\"he has filled the hungry with good things, and sent the rich away empty.\"",
                    "ref": "Luke 1:53"
                },
                {
                    "text": "\"He has helped his servant Israel, in remembrance of his mercy,\"",
                    "ref": "Luke 1:54"
                },
                {
                    "text": "\"according to the promise he made to our ancestors, to Abraham and to his descendants forever.\"",
                    "ref": "Luke 1:55"
                }
            ]
        },
        {
            "title": "Coronation of Mary",
            "focus": "Trust in Mary's intercession.",
            "verses": [
                {
                    "text": "A great portent appeared in heaven: a woman clothed with the sun, with the moon under her feet, and on her head a crown of twelve stars.",
                    "ref": "Revelation 12:1"
                },
                {
                    "text": "She was pregnant and was crying out in birth pangs, in the agony of giving birth.",
                    "ref": "Revelation 12:2"
                },
                {
                    "text": "Then another portent appeared in heaven: a great red dragon, with seven heads and ten horns, and seven diadems on his heads.",
                    "ref": "Revelation 12:3"
                },
                {
                    "text": "His tail swept down a third of the stars of heaven and threw them to the earth.",
                    "ref": "Revelation 12:4"
                },
                {
                    "text": "And she gave birth to a son, a male child, who is to rule all the nations with a rod of iron. But her child was snatched away to God and to his throne;",
                    "ref": "Revelation 12:5"
                },
                {
                    "text": "And war broke out in heaven; Michael and his angels fought against the dragon. The dragon and his angels fought back,",
                    "ref": "Revelation 12:7"
                },
                {
                    "text": "but they were defeated, and there was no longer any place for them in heaven.",
                    "ref": "Revelation 12:8"
                },
                {
                    "text": "The great dragon was thrown down, that ancient serpent, who is called the Devil and Satan, the deceiver of the whole world—he was thrown down to the earth, and his angels were thrown down with him.",
                    "ref": "Revelation 12:9"
                },
                {
                    "text": "Then I heard a loud voice in heaven, proclaiming, 'Now have come the salvation and the power and the kingdom of our God and the authority of his Messiah...'",
                    "ref": "Revelation 12:10"
                },
                {
                    "text": "But they have conquered him by the blood of the Lamb and by the word of their testimony, for they did not cling to life even in the face of death.",
                    "ref": "Revelation 12:11"
                }
            ]
        }
    ]
};

export const generateRosarySteps = (mysterySet: string): RosaryStep[] => {
    const steps: RosaryStep[] = [];
    const mysteries = ROSARY_DATA[mysterySet] || [];

    // 1. Intro
    steps.push({ type: 'opening', prayerId: 'sign-of-the-cross' });
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
            content: mystery.focus
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
    steps.push({ type: 'opening', prayerId: 'sign-of-the-cross' });

    return steps;
};
