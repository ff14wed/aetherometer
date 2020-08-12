package testassets

const ActionCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65
#,Name,,Icon,ActionCategory,,Animation{Start},VFX,Animation{End},ActionTimeline{Hit},,ClassJob,BehaviourType,ClassJobLevel,IsRoleAction,Range,CanTargetSelf,CanTargetParty,CanTargetFriendly,CanTargetHostile,,,TargetArea,,,,CanTargetDead,,CastType,EffectRange,XAxisModifier,,PrimaryCost{Type},PrimaryCost{Value},SecondaryCost{Type},SecondaryCost{Value},Action{Combo},PreservesCombo,Cast<100ms>,Recast<100ms>,CooldownGroup,,MaxCharges,AttackType,Aspect,ActionProcStatus,,Status{GainSelf},UnlockLink,ClassJobCategory,,,AffectsPosition,Omen,IsPvP,,,,,,,,,,,,IsPlayerAction
int32,str,bit&01,Image,ActionCategory,byte,ActionCastTimeline,ActionCastVFX,ActionTimeline,ActionTimeline,byte,ClassJob,byte,byte,bit&02,sbyte,bit&04,bit&08,bit&10,bit&20,bit&40,bit&80,bit&01,bit&02,bit&04,sbyte,bit&08,bit&10,byte,byte,byte,bit&20,byte,uint16,byte,Row,Action,bit&40,uint16,uint16,byte,byte,byte,AttackType,byte,ActionProcStatus,byte,Status,Row,ClassJobCategory,byte,bit&80,bit&01,Omen,bit&02,bit&04,bit&08,bit&10,bit&20,bit&40,bit&80,bit&01,bit&02,byte,bit&04,bit&08,bit&10
0,"",False,405,0,0,0,0,0,0,0,-1,0,0,False,0,False,False,False,False,False,False,False,False,False,0,True,True,0,0,0,False,0,0,0,0,0,False,0,0,0,0,0,0,0,0,0,0,0,0,0,False,False,0,False,False,True,False,True,True,False,False,True,0,False,False,False
2,"Interaction",False,405,8,0,1,0,0,1875,0,0,0,0,False,3,False,False,False,False,False,False,False,False,False,0,True,True,1,0,0,False,0,0,0,0,0,False,50,0,0,0,0,0,0,0,0,0,0,1,0,False,False,0,False,False,True,False,True,True,False,False,True,0,False,False,False
3,"Sprint",False,405,10,0,0,0,368,1875,0,0,0,0,False,0,True,False,False,False,False,False,False,False,False,0,True,True,1,0,0,False,0,0,0,0,0,True,0,600,56,0,0,0,0,0,0,0,0,1,0,False,False,0,False,False,True,False,True,True,False,False,True,0,False,True,False
4,"Mount",False,405,5,0,2,0,165,1875,0,0,0,0,False,0,True,False,False,False,False,False,False,False,False,0,True,True,1,0,0,False,0,0,0,0,0,False,10,0,0,0,0,0,0,0,0,0,0,1,0,False,False,0,False,False,False,False,True,True,False,False,True,0,False,True,False
5,"Teleport",False,111,10,0,1,1,164,1875,0,0,0,0,False,0,True,False,False,False,False,False,False,False,False,0,True,True,1,0,0,False,0,0,0,0,0,False,50,0,0,0,0,0,0,0,0,0,4,1,0,False,False,0,False,False,True,False,True,True,False,False,True,0,True,True,False
7,"attack",False,405,1,0,0,0,-1,1875,0,0,0,0,False,-1,False,False,False,True,False,False,False,False,False,0,True,True,1,0,0,True,0,0,0,0,0,True,0,0,0,0,0,-1,7,0,1,0,0,1,0,False,False,0,False,False,True,False,True,True,False,False,True,0,False,False,False
9,"Fast Blade",False,158,3,0,0,0,310,1875,0,1,0,1,False,-1,False,False,False,True,False,False,False,False,False,0,True,True,1,0,0,True,0,0,0,0,0,False,0,25,58,0,0,-1,7,0,1,0,0,38,2,False,False,0,False,False,True,False,True,True,False,False,True,0,False,False,True
11,"Savage Blade",False,157,3,0,0,0,311,1875,0,-1,0,0,False,-1,False,False,False,True,False,False,False,False,False,0,True,True,1,0,0,True,0,0,0,0,9,False,0,25,58,0,0,-1,7,0,1,0,0,0,2,False,False,0,False,False,True,False,True,True,False,False,True,0,False,False,False
26,"Sword Oath",False,2504,2,0,29,4,442,1875,0,-1,0,0,False,0,True,False,False,False,False,False,False,False,False,0,True,True,1,0,0,False,3,10,0,0,0,True,0,25,58,0,0,0,0,0,2,381,0,0,1,True,False,0,False,False,True,False,True,True,False,False,True,0,False,False,False
50,"Unchained",False,2554,4,0,0,0,386,1875,0,-1,0,0,False,0,True,False,False,False,False,False,False,False,False,0,True,True,1,0,0,False,0,0,32,91,0,True,0,900,8,0,0,0,0,0,2,0,0,0,1,False,False,0,False,False,True,False,True,True,False,False,True,0,False,False,False
102,"Flaming Arrow",False,368,4,0,0,0,416,1875,0,-1,0,0,False,-1,False,False,False,False,False,False,True,False,False,0,True,True,7,5,0,False,0,0,0,0,0,True,0,600,0,0,0,-1,1,0,1,0,0,0,1,False,False,0,False,False,True,False,True,True,False,False,True,0,False,False,False
203,"Skyshard",False,103,9,1,52,22,2251,1883,0,0,0,0,False,25,False,False,False,False,False,False,True,False,False,0,True,True,2,8,0,True,11,0,0,0,0,False,20,0,0,0,0,8,7,0,1,0,0,1,0,False,False,1,False,False,True,False,True,True,True,False,True,0,False,False,False
4238,"Big Shot",False,103,9,0,54,31,3925,4012,0,0,0,0,False,30,False,False,False,True,False,False,False,False,False,0,True,True,4,30,4,True,11,0,0,0,0,False,20,0,0,0,0,8,7,0,1,0,0,1,0,False,False,2,False,False,True,False,True,True,True,False,True,0,False,False,False
`

const OmenCSV = `
key,0,1,2,3,4,5
#,Path,PathAlly,Type,RestrictYScale,LargeScale,
int32,str,str,byte,bit&01,bit&02,sbyte
0,"","",0,True,False,0
1,"general_1bf","general_1bpf",0,True,False,0
2,"general02f","general02_pf",0,True,False,0
`

const CraftActionCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19
#,Name,Description,Animation{Start},Animation{End},Icon,ClassJob,ClassJobCategory,ClassJobLevel,QuestRequirement,Specialist,,Cost,CRP,BSM,ARM,GSM,LTW,WVR,ALC,CUL
int32,str,str,ActionTimeline,ActionTimeline,Image,ClassJob,ClassJobCategory,byte,Quest,bit&01,uint16,byte,CraftAction,CraftAction,CraftAction,CraftAction,CraftAction,CraftAction,CraftAction,CraftAction
100000,"","",0,0,0,0,0,0,0,False,0,0,0,0,0,0,0,0,0,0
100001,"Basic Synthesis","Increases progress.
<UIForeground>F201F8</UIForeground><UIGlow>F201F9</UIGlow>Efficiency:<UIGlow>01</UIGlow><UIForeground>01</UIForeground> <If(Equal(PlayerParameter(68),8))><If(GreaterThanOrEqualTo(PlayerParameter(69),31))>120<Else/>100</If><Else/>100</If>%
<UIForeground>F201F8</UIForeground><UIGlow>F201F9</UIGlow>Success Rate:<UIGlow>01</UIGlow><UIForeground>01</UIForeground> 100%",239,246,1501,8,9,1,0,False,0,0,100001,100015,100030,100075,100045,100060,100090,100105
100002,"Basic Touch","Increases quality.
<UIForeground>F201F8</UIForeground><UIGlow>F201F9</UIGlow>Efficiency:<UIGlow>01</UIGlow><UIForeground>01</UIForeground> 100%
<UIForeground>F201F8</UIForeground><UIGlow>F201F9</UIGlow>Success Rate:<UIGlow>01</UIGlow><UIForeground>01</UIForeground> 100%",240,247,1502,8,9,5,0,False,0,18,100002,100016,100031,100076,100046,100061,100091,100106
`

const BNPCBaseCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20
#,Behavior,Battalion,LinkRace,Rank,Scale,ModelChara,BNpcCustomize,NpcEquip,Special,SEPack,,ArrayEventHandler,BNpcParts,,IsTargetLine,IsDisplayLevel,,,,,
int32,Behavior,Battalion,LinkRace,byte,single,ModelChara,BNpcCustomize,NpcEquip,uint16,byte,bit&01,ArrayEventHandler,BNpcParts,bit&02,bit&04,bit&08,bit&10,bit&20,byte,byte,byte
0,0,4,0,0,1,0,0,0,0,0,False,0,0,True,True,True,True,False,0,0,0
1,0,4,0,0,1,0,2,366,0,0,False,0,0,True,True,True,True,False,0,0,0
2,0,4,0,0,1,96,0,0,0,0,False,851976,0,True,True,True,True,False,0,0,0
3,0,4,0,0,1.2,61,0,0,0,0,False,852023,0,True,True,True,True,False,0,0,0
`

const BNPCNameCSV = `
key,0,1,2,3,4,5,6,7
#,Singular,Adjective,Plural,PossessivePronoun,StartsWithVowel,,Pronoun,Article
int32,str,sbyte,str,sbyte,sbyte,sbyte,sbyte,sbyte
0,"",0,"",0,0,1,0,0
1,"",0,"",0,0,0,0,0
2,"ruins runner",0,"ruins runners",0,0,1,0,0
3,"antelope doe",0,"antelope does",0,1,1,0,0
`

const ModelCharaCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16
#,Type,Model,Base,Variant,SEPack,,,PapVariation,,,,,,,,,
int32,byte,uint16,byte,byte,uint16,byte,bit&01,bit&02,byte,sbyte,bit&04,bit&08,bit&10,bit&20,bit&40,byte,bit&80
878,3,8094,1,1,3419,1,False,False,0,0,False,False,False,False,False,100,False
879,3,8095,1,1,3419,1,False,False,0,0,False,False,False,False,False,100,False
880,3,8096,1,1,3419,1,False,False,0,0,False,False,False,False,False,100,False
881,3,8097,1,1,3419,1,False,False,0,0,False,False,False,False,False,100,False
882,3,8098,1,1,3419,1,False,False,0,0,False,False,False,False,False,100,False
883,3,8099,1,1,3479,1,False,False,0,0,False,False,False,False,False,100,False
`

const ModelSkeletonCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16
#,Radius,Height,VFXScale,,,,,,,,,FloatHeight,FloatDown,FloatUp,,MotionBlendType,LoopFlySE
int32,single,single,single,uint16,uint16,uint16,uint16,uint16,uint16,uint16,uint16,single,single,uint16,byte,bit&01,byte
8094,0.2,0.3,0.5,100,300,66,300,66,300,0,0,3,3,0,0,False,1
8095,0.2,0.3,0.5,100,300,66,300,66,300,0,0,3,3,0,0,False,1
8096,0.2,0.3,0.5,100,300,66,300,66,300,0,0,3,3,0,0,False,1
8097,0.2,0.3,0.5,100,300,66,300,66,300,0,0,3,3,0,0,False,1
8098,0.2,0.3,0.5,100,300,66,300,66,300,0,0,3,3,0,0,False,1
`

const MapCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17
#,MapCondition,PriorityCategoryUI,PriorityUI,MapIndex,Hierarchy,MapMarkerRange,Id,SizeFactor,Offset{X},Offset{Y},PlaceName{Region},PlaceName,PlaceName{Sub},DiscoveryIndex,DiscoveryFlag,TerritoryType,DiscoveryArrayByte,IsEvent
int32,MapCondition,byte,byte,sbyte,byte,uint16,str,uint16,int16,int16,PlaceName,PlaceName,PlaceName,int16,uint32,TerritoryType,bit&01,bit&02
0,0,0,0,0,0,0,"",100,0,0,0,0,0,-1,0,0,True,False
1,0,0,0,0,1,0,"default/00",100,0,0,2405,21,0,-1,0,1,True,False
2,0,2,2,0,1,3,"f1t1/00",200,0,0,23,52,0,-1,0,132,True,False
3,0,2,3,0,1,4,"f1t2/00",200,0,0,23,53,0,-1,0,133,True,False
14,0,3,3,1,1,60,"w1t2/01",200,0,0,24,41,373,-1,0,131,True,False
73,0,3,3,2,1,62,"w1t2/02",200,0,0,24,41,698,-1,0,131,True,False
178,0,0,0,0,1,143,"w1b4/00",200,-448,0,24,1409,0,-1,0,196,True,False
33,0,0,0,0,1,0,"s1fa/00",400,0,0,22,359,0,-1,0,206,True,False
403,0,0,0,0,1,0,"s1fa/00",400,0,0,22,359,19,-1,0,293,True,False
`

const PlaceNameCSV = `
key,0,1,2,3,4,5,6,7,8,9
#,Name,,Name{NoArticle},,,,,,,
int32,str,sbyte,str,sbyte,sbyte,sbyte,sbyte,sbyte,str,byte
0,"",0,"",0,0,1,0,0,"",0
19,"",0,"",0,0,1,0,0,"",0
21,"Eorzea",1,"Eorzea",0,0,1,0,1,"",0
41,"Ul'dah - Steps of Thal",2,"Ul'dah - Steps of Thal",0,0,1,0,1,"",0
52,"New Gridania",1,"New Gridania",0,0,1,0,1,"",0
53,"Old Gridania",1,"Old Gridania",0,0,1,0,1,"",0
359,"The Navel",1,"Navel",0,0,1,0,0,"",0
373,"Merchant Strip",2,"Merchant Strip",0,0,1,0,0,"",0
698,"Hustings Strip",2,"Hustings Strip",0,0,1,0,0,"",0
1409,"The Burning Heart",0,"Burning Heart",0,0,1,0,0,"",0
`

const TerritoryTypeCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39
#,Name,Bg,BattalionMode,PlaceName{Region},PlaceName{Zone},PlaceName,Map,LoadingImage,ExclusiveType,TerritoryIntendedUse,,,WeatherRate,,,PCSearch,Stealth,Mount,,BGM,PlaceName{Region}Icon,PlaceNameIcon,ArrayEventHandler,QuestBattle,Aetheryte,FixedTime,Resident,AchievementIndex,IsPvpZone,ExVersion,,,AddedIn{5-3},MountSpeed,,,,,,
int32,str,str,byte,PlaceName,PlaceName,PlaceName,Map,LoadingImage,byte,byte,uint16,bit&01,byte,bit&02,byte,bit&04,bit&08,bit&10,bit&20,BGM,Image,Image,ArrayEventHandler,QuestBattle,Aetheryte,int32,uint16,sbyte,bit&40,ExVersion,byte,byte,byte,MountSpeed,bit&80,bit&01,byte,bit&02,bit&04,bit&08
1,"","",0,0,0,0,0,0,0,0,0,False,0,False,0,False,False,False,False,0,0,0,0,0,0,0,0,0,False,0,0,0,0,0,False,False,0,False,False,False
128,"s1t1","ffxiv/sea_s1/twn/s1t1/level/s1t1",1,22,500,28,11,2,0,0,0,False,14,True,0,True,False,False,False,1020,122007,123002,852085,0,8,-1,0,0,False,0,0,0,0,0,False,False,0,False,False,False
129,"s1t2","ffxiv/sea_s1/twn/s1t2/level/s1t2",1,22,500,29,12,2,0,0,0,False,15,True,0,True,False,False,False,1020,122007,123003,852088,0,8,-1,0,1,False,0,0,0,0,0,False,False,0,False,False,False
130,"w1t1","ffxiv/wil_w1/twn/w1t1/level/w1t1",1,24,504,40,13,4,0,0,0,False,7,True,0,True,False,False,False,1035,122008,123102,852086,0,9,-1,0,2,False,0,0,0,0,0,False,False,0,False,False,False
131,"w1t2","ffxiv/wil_w1/twn/w1t2/level/w1t2",1,24,504,41,14,4,0,0,0,False,8,True,0,True,False,False,False,1035,122008,123103,853381,0,9,-1,0,3,False,0,0,0,0,0,False,False,0,False,False,False
132,"f1t1","ffxiv/fst_f1/twn/f1t1/level/f1t1",1,23,506,52,2,3,0,0,0,False,1,True,0,True,False,False,False,1003,122009,123202,852087,0,2,-1,2,4,False,0,0,0,0,0,False,False,0,False,False,False
133,"f1t2","ffxiv/fst_f1/twn/f1t2/level/f1t2",1,23,506,53,3,3,0,0,0,False,2,True,0,True,False,False,False,1003,122009,123203,852103,0,2,-1,0,5,False,0,0,0,0,0,False,False,0,False,False,False
196,"w1b4","ffxiv/wil_w1/bah/w1b4/level/w1b4",1,24,505,1409,178,4,2,17,110,False,44,False,0,True,False,False,False,1001,122010,124534,0,0,0,-1,63,-1,False,0,0,0,0,0,False,False,0,False,False,False
206,"s1fa","ffxiv/sea_s1/fld/s1fa/level/s1fa",1,22,502,359,33,2,2,10,57,False,23,False,0,True,False,False,False,1001,-1,-1,0,0,0,-1,6,-1,False,0,0,0,0,0,False,False,0,False,False,False
293,"s1fa_2","ffxiv/sea_s1/fld/s1fa/level/s1fa",1,22,502,359,403,2,2,10,60,False,23,False,0,True,False,False,False,1001,122001,124009,0,0,0,-1,6,-1,False,0,0,0,0,0,False,False,0,False,False,False
296,"s1fa_3","ffxiv/sea_s1/fld/s1fa/level/s1fa",1,22,502,359,403,2,2,10,64,False,23,False,0,True,False,False,False,1001,122001,124010,0,0,0,-1,16,-1,False,0,0,0,0,0,False,False,0,False,False,False
`

const StatusCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29
#,Name,Description,Icon,MaxStacks,,Category,HitEffect,VFX,LockMovement,,LockActions,LockControl,Transfiguration,,CanDispel,InflictedByActor,IsPermanent,PartyListPriority,,,,,,,Log,IsFcBuff,Invisibility,,,
int32,str,str,Image,byte,byte,byte,StatusHitEffect,StatusLoopVFX,bit&01,bit&02,bit&04,bit&08,bit&10,bit&20,bit&40,bit&80,bit&01,byte,bit&02,bit&04,bit&08,int16,byte,bit&10,uint16,bit&20,bit&40,byte,byte,bit&80
0,"","",0,0,0,0,0,0,False,False,False,False,False,False,False,False,False,0,False,False,False,0,0,False,0,False,False,0,0,False
1,"Petrification","Stone-like rigidity is preventing the execution of actions.",15001,0,1,2,6,1,True,False,True,False,False,False,False,False,False,100,False,False,False,0,0,False,0,False,False,0,32,True
2,"Stun","Unable to execute actions.",15004,0,1,2,3,2,True,False,True,False,False,False,True,False,False,100,False,False,False,0,0,False,0,False,False,0,0,False
`

const ClassJobCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45
#,Name,Abbreviation,,ClassJobCategory,ExpArrayIndex,BattleClassIndex,,JobIndex,,Modifier{HitPoints},Modifier{ManaPoints},Modifier{Strength},Modifier{Vitality},Modifier{Dexterity},Modifier{Intelligence},Modifier{Mind},Modifier{Piety},,,,,,,,,,ClassJob{Parent},Name{English},Item{StartingWeapon},,Role,StartingTown,MonsterNote,PrimaryStat,LimitBreak1,LimitBreak2,LimitBreak3,UIPriority,Item{SoulCrystal},UnlockQuest,RelicQuest,Prerequisite,StartingLevel,PartyBonus,IsLimitedJob,CanQueueForDuty
int32,str,str,str,ClassJobCategory,sbyte,sbyte,byte,byte,sbyte,uint16,uint16,uint16,uint16,uint16,uint16,uint16,uint16,uint16,uint16,uint16,uint16,uint16,uint16,byte,byte,byte,ClassJob,str,Item,int32,byte,Town,MonsterNote,byte,Action,Action,Action,byte,Item,Quest,Quest,Quest,byte,byte,bit&01,bit&01
0,"adventurer","ADV","",30,-1,-1,0,0,-1,100,100,100,100,100,100,100,100,100,100,100,100,100,100,0,0,0,0,"Adventurer",0,0,0,0,127,0,0,0,0,0,0,0,0,0,1,0,False,False
1,"gladiator","GLA","剣",30,1,0,1,0,-1,110,100,95,100,90,50,95,100,100,100,100,100,100,100,0,0,0,1,"Gladiator",1601,0,1,3,0,1,197,198,199,2,0,0,0,0,1,1,False,True
2,"pugilist","PGL","格",30,0,1,2,0,-1,105,100,100,95,100,45,85,100,100,100,100,100,100,100,0,0,0,2,"Pugilist",1680,0,2,3,1,1,200,201,202,22,0,0,0,0,1,3,False,True
`

const RecipeCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45
#,Number,CraftType,RecipeLevelTable,Item{Result},Amount{Result},Item{Ingredient}[0],Amount{Ingredient}[0],Item{Ingredient}[1],Amount{Ingredient}[1],Item{Ingredient}[2],Amount{Ingredient}[2],Item{Ingredient}[3],Amount{Ingredient}[3],Item{Ingredient}[4],Amount{Ingredient}[4],Item{Ingredient}[5],Amount{Ingredient}[5],Item{Ingredient}[6],Amount{Ingredient}[6],Item{Ingredient}[7],Amount{Ingredient}[7],Item{Ingredient}[8],Amount{Ingredient}[8],Item{Ingredient}[9],Amount{Ingredient}[9],,IsSecondary,MaterialQualityFactor,DifficultyFactor,QualityFactor,DurabilityFactor,,RequiredCraftsmanship,RequiredControl,QuickSynthCraftsmanship,QuickSynthControl,SecretRecipeBook,,CanQuickSynth,CanHq,ExpRewarded,Status{Required},Item{Required},IsSpecializationRequired,IsExpert,PatchNumber
int32,int32,CraftType,RecipeLevelTable,Item,byte,Item,byte,Item,byte,Item,byte,Item,byte,Item,byte,Item,byte,Item,byte,Item,byte,Item,byte,Item,byte,uint16,bit&01,byte,uint16,uint16,uint16,uint16,uint16,uint16,uint16,uint16,SecretRecipeBook,uint32,bit&02,bit&04,bit&08,Status,Item,bit&10,bit&20,uint16
1,10001,1,1,5056,1,5106,2,5107,1,0,0,0,0,0,0,0,0,0,0,0,0,2,1,-1,0,0,False,50,50,100,67,0,0,0,0,0,0,0,True,True,True,0,0,False,False,0
33067,14536,5,320,23002,1,23375,2,22493,2,19988,1,0,0,0,0,0,0,0,0,0,0,18,2,16,2,1045,False,50,100,100,100,0,1320,1220,1500,0,61,0,True,True,True,0,0,False,False,430
33068,13223,4,320,23374,1,23376,2,22430,1,19946,1,19988,1,0,0,0,0,0,0,0,0,17,2,16,2,1044,False,50,100,100,100,0,1320,1220,1500,0,60,0,True,True,True,0,0,False,False,430
33073,10833,1,380,23768,1,24250,2,19943,1,24258,2,0,0,0,0,0,0,0,0,0,0,14,2,17,2,1041,False,50,100,100,100,0,1650,1600,0,0,57,0,False,True,True,0,0,False,False,440
33074,10837,1,380,23769,1,24250,3,24251,1,19943,1,24256,3,0,0,0,0,0,0,0,0,14,2,17,2,1041,False,50,100,100,100,0,1650,1600,0,0,57,0,False,True,True,0,0,False,False,440
`

const RecipeLevelTableCSV = `
key,0,1,2,3,4,5,6
#,ClassJobLevel,Stars,SuggestedCraftsmanship,SuggestedControl,Difficulty,Quality,Durability
int32,byte,byte,uint16,uint16,uint16,uint32,uint16
320,70,2,1320,1220,2657,12511,70
380,70,4,1650,1600,2900,14267,70
`

const ItemCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90
#,Singular,Adjective,Plural,PossessivePronoun,StartsWithVowel,,Pronoun,Article,Description,Name,Icon,Level{Item},Rarity,FilterGroup,AdditionalData,ItemUICategory,ItemSearchCategory,EquipSlotCategory,,AddedIn{5-3},StackSize,IsUnique,IsUntradable,IsIndisposable,Lot,Price{Mid},Price{Low},CanBeHq,IsDyeable,IsCrestWorthy,ItemAction,,Cooldown<s>,ClassJob{Repair},Item{Repair},Item{Glamour},,IsCollectable,AlwaysCollectable,AetherialReduce,Level{Equip},,EquipRestriction,ClassJobCategory,GrandCompany,ItemSeries,BaseParamModifier,Model{Main},Model{Sub},ClassJob{Use},,Damage{Phys},Damage{Mag},Delay<ms>,,BlockRate,Block,Defense{Phys},Defense{Mag},BaseParam[0],BaseParamValue[0],BaseParam[1],BaseParamValue[1],BaseParam[2],BaseParamValue[2],BaseParam[3],BaseParamValue[3],BaseParam[4],BaseParamValue[4],BaseParam[5],BaseParamValue[5],ItemSpecialBonus,ItemSpecialBonus{Param},BaseParam{Special}[0],BaseParamValue{Special}[0],BaseParam{Special}[1],BaseParamValue{Special}[1],BaseParam{Special}[2],BaseParamValue{Special}[2],BaseParam{Special}[3],BaseParamValue{Special}[3],BaseParam{Special}[4],BaseParamValue{Special}[4],BaseParam{Special}[5],BaseParamValue{Special}[5],MaterializeType,MateriaSlotCount,IsAdvancedMeldingPermitted,IsPvP,,IsGlamourous
int32,str,sbyte,str,sbyte,sbyte,sbyte,sbyte,sbyte,str,str,Image,ItemLevel,byte,byte,Row,ItemUICategory,ItemSearchCategory,EquipSlotCategory,byte,uint16,uint32,bit&01,bit&02,bit&04,bit&08,uint32,uint32,bit&10,bit&20,bit&40,ItemAction,byte,uint16,ClassJob,Item,Item,uint16,bit&80,bit&01,uint16,byte,byte,byte,ClassJobCategory,GrandCompany,ItemSeries,byte,int64,int64,ClassJob,byte,uint16,uint16,uint16,byte,uint16,uint16,uint16,uint16,BaseParam,int16,BaseParam,int16,BaseParam,int16,BaseParam,int16,BaseParam,int16,BaseParam,int16,ItemSpecialBonus,byte,BaseParam,int16,BaseParam,int16,BaseParam,int16,BaseParam,int16,BaseParam,int16,BaseParam,int16,byte,byte,bit&01,bit&02,byte,bit&04
5056,"bronze ingot",0,"bronze ingots",0,0,1,0,0,"An ingot of smelted bronze.","Bronze Ingot",20803,1,1,12,0,49,48,0,16,2000,999,False,False,False,False,9,1,True,False,False,0,2,0,0,0,0,0,False,False,0,1,0,0,0,0,0,0,"0, 0, 0, 0","0, 0, 0, 0",0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,False,False,0,False
23374,"pair of quaintrelle's dress shoes",0,"pairs of quaintrelle's dress shoes",0,0,1,0,0,"Fits: All ♀","Quaintrelle's Dress Shoes",46071,1,1,4,0,38,37,8,5,64550,1,False,False,False,False,94,1,True,True,False,0,2,0,12,5594,21800,2997,False,False,0,1,0,3,1,0,0,0,"6070, 1, 0, 0","0, 0, 0, 0",0,0,0,0,0,0,0,0,4,9,0,0,0,0,0,0,0,0,0,0,0,0,1,0,21,2,24,2,0,0,0,0,0,0,0,0,1,0,False,False,0,True
23768,"rakshasa blade",0,"rakshasa blades",0,0,1,0,0,"","Rakshasa Blade",30592,380,2,1,0,2,10,1,5,2320,1,False,False,False,False,71640,414,True,True,False,0,2,0,9,17837,21800,428,False,False,0,70,0,1,38,0,0,1,"201, 64, 3, 0","0, 0, 0, 0",1,1,95,64,2080,3,0,0,0,0,1,233,3,259,45,208,44,146,0,0,0,0,1,0,12,10,13,7,1,26,3,28,45,23,44,16,3,2,True,False,0,True
23769,"pair of rakshasa knuckles",0,"pairs of rakshasa knuckles",0,0,1,0,0,"","Rakshasa Knuckles",31158,380,2,1,0,1,9,13,5,16320,1,False,False,False,False,107460,621,True,True,False,0,2,0,9,17837,21800,425,False,False,0,70,0,1,41,0,0,3,"323, 29, 1, 0","373, 29, 1, 0",2,3,95,64,2560,3,0,0,0,0,1,327,3,362,45,292,22,204,0,0,0,0,1,0,12,10,13,7,1,36,3,40,45,32,22,23,5,2,True,False,0,True
`

const WorldCSV = `
key,0,1,2,3
#,Name,UserType,DataCenter,IsPublic
int32,str,byte,WorldDCGroupType,bit&01
0,"crossworld",0,0,False
1,"reserved1",0,1,False
2,"c-contents",0,0,False
3,"c-whiteae",0,1,False
4,"c-baudinii",0,0,False
5,"c-contents2",0,2,False
`
